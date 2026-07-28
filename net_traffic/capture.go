package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

/*
run
sudo ./capture -i ens192 -interval 2s	#2초 주기로 ens192에서 수집한 패킷정보
sudo ./capture -i ens192 -interval 2s -filter "tcp port 9083" # 위 패킷정보에 port 필터링
*/

// Capturer 캡처 수행자: 인터페이스, handle, 통계(원자)
type Capturer struct {
	iface     string
	snaplen   int32
	promisc   bool
	timeout   time.Duration
	filter    string
	handle    *pcap.Handle
	packetCnt uint64 // atomic
	byteCnt   uint64 // atomic
}

// NewCapturer 생성자
func NewCapturer(iface string, snaplen int32, promisc bool, timeout time.Duration, filter string) *Capturer {
	return &Capturer{
		iface:   iface,
		snaplen: snaplen,
		promisc: promisc,
		timeout: timeout,
		filter:  filter,
	}
}

// Start 캡처를 시작한다. ctx가 취소되면 리턴.
func (c *Capturer) Start(ctx context.Context) error {
	handle, err := pcap.OpenLive(c.iface, int32(c.snaplen), c.promisc, c.timeout)
	if err != nil {
		return fmt.Errorf("pcap.OpenLive: %w", err)
	}
	c.handle = handle

	// optional filter
	if c.filter != "" {
		if err := handle.SetBPFFilter(c.filter); err != nil {
			// close handle on error
			handle.Close()
			return fmt.Errorf("SetBPFFilter(%s): %w", c.filter, err)
		}
	}

	src := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := src.Packets()

	// read loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case pkt, ok := <-packets:
				if !ok {
					// packet source closed
					return
				}
				if pkt == nil {
					continue
				}
				data := pkt.Data()
				atomic.AddUint64(&c.packetCnt, 1)
				atomic.AddUint64(&c.byteCnt, uint64(len(data)))
			}
		}
	}()

	return nil
}

// Close handle 닫기
func (c *Capturer) Close() {
	if c.handle != nil {
		c.handle.Close()
		c.handle = nil
	}
}

// SnapshotStats와 리셋 (원자적으로 읽고 0으로)
func (c *Capturer) SnapshotAndReset() (pkts uint64, bytes uint64) {
	pkts = atomic.SwapUint64(&c.packetCnt, 0)
	bytes = atomic.SwapUint64(&c.byteCnt, 0)
	return
}

func main() {
	iface := flag.String("i", "ens192", "network interface to capture (e.g. ens192)")
	interval := flag.Duration("interval", 1*time.Second, "report interval (e.g. 1s, 500ms)")
	snaplen := flag.Int("snaplen", 65535, "snapshot length for pcap")
	promisc := flag.Bool("promisc", false, "enable promiscuous mode")
	filter := flag.String("filter", "", "BPF filter (optional)")
	flag.Parse()

	// prepare capturer
	c := NewCapturer(*iface, int32(*snaplen), *promisc, pcap.BlockForever, *filter)

	// context + signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	if err := c.Start(ctx); err != nil {
		log.Fatalf("failed to start capturer: %v", err)
	}
	defer c.Close()

	log.Printf("started capture on interface=%s interval=%s filter=%q", *iface, interval.String(), *filter)

	// reporting loop
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			pkts, bytes := c.SnapshotAndReset()
			pktRate := float64(pkts) / interval.Seconds()
			byteRate := float64(bytes) / interval.Seconds()
			log.Printf("interval %s: packets=%d (pkt/s=%.2f), bytes=%d (B/s=%.2f)", interval.String(), pkts, pktRate, bytes, byteRate)
		case s := <-sigCh:
			log.Printf("signal %v received, shutting down...", s)
			cancel()
			// wait a short moment to ensure goroutine exits and remaining packets are read
			time.Sleep(200 * time.Millisecond)
			break loop
		}
	}

	// final stats flush
	pkts, bytes := c.SnapshotAndReset()
	if pkts > 0 || bytes > 0 {
		log.Printf("final stats: packets=%d, bytes=%d", pkts, bytes)
	}
	log.Println("exited")
}
