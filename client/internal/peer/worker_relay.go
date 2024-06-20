package peer

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	relayClient "github.com/netbirdio/netbird/relay/client"
)

type RelayConnInfo struct {
	relayedConn     net.Conn
	rosenpassPubKey []byte
	rosenpassAddr   string
}

type WorkerRelayCallbacks struct {
	OnConnReady     func(RelayConnInfo)
	OnStatusChanged func(ConnStatus)
	DoHandshake     func() (*OfferAnswer, error)
}

type WorkerRelay struct {
	ctx          context.Context
	log          *log.Entry
	relayManager *relayClient.Manager
	config       ConnConfig
	conn         WorkerRelayCallbacks
}

func NewWorkerRelay(ctx context.Context, log *log.Entry, relayManager *relayClient.Manager, config ConnConfig, callbacks WorkerRelayCallbacks) *WorkerRelay {
	return &WorkerRelay{
		ctx:          ctx,
		log:          log,
		relayManager: relayManager,
		config:       config,
		conn:         callbacks,
	}
}

// SetupRelayConnection todo: this function is not completed. Make no sense to put it in a for loop because we are not waiting for any event
func (w *WorkerRelay) SetupRelayConnection() {
	for {
		if !w.waitForReconnectTry() {
			return
		}

		w.log.Debugf("trying to establish Relay connection with peer %s", w.config.Key)

		remoteOfferAnswer, err := w.conn.DoHandshake()
		if err != nil {
			if errors.Is(err, ErrSignalIsNotReady) {
				w.log.Infof("signal client isn't ready, skipping connection attempt")
			}
			w.log.Errorf("%s", err)
			continue
		}

		if !w.isRelaySupported(remoteOfferAnswer) {
			w.log.Infof("Relay is not supported by remote peer")
			// todo should we retry?
			// if the remote peer doesn't support relay make no sense to retry infinity
			// but if the remote peer supports relay just the connection is lost we should retry
			continue
		}

		// the relayManager will return with error in case if the connection has lost with relay server
		currentRelayAddress, err := w.relayManager.RelayAddress()
		if err != nil {
			w.log.Infof("local Relay connection is lost, skipping connection attempt")
			continue
		}

		srv := w.preferredRelayServer(currentRelayAddress.String(), remoteOfferAnswer.RelaySrvAddress)
		relayedConn, err := w.relayManager.OpenConn(srv, w.config.Key)
		if err != nil {
			w.log.Infof("failed to open relay connection: %s", err)
			continue
		}

		go w.conn.OnConnReady(RelayConnInfo{
			relayedConn:     relayedConn,
			rosenpassPubKey: remoteOfferAnswer.RosenpassPubKey,
			rosenpassAddr:   remoteOfferAnswer.RosenpassAddr,
		})

		<-w.ctx.Done()
	}
}

func (w *WorkerRelay) RelayAddress() (net.Addr, error) {
	return w.relayManager.RelayAddress()
}

func (w *WorkerRelay) isRelaySupported(answer *OfferAnswer) bool {
	if !w.relayManager.HasRelayAddress() {
		return false
	}
	return answer.RelaySrvAddress != ""
}

func (w *WorkerRelay) preferredRelayServer(myRelayAddress, remoteRelayAddress string) string {
	if w.config.LocalKey > w.config.Key {
		return myRelayAddress
	}
	return remoteRelayAddress
}

func (w *WorkerRelay) RelayIsSupportedLocally() bool {
	return w.relayManager.HasRelayAddress()
}

// waitForReconnectTry waits for a random duration before trying to reconnect
func (w *WorkerRelay) waitForReconnectTry() bool {
	minWait := 500
	maxWait := 2000
	duration := time.Duration(rand.Intn(maxWait-minWait)+minWait) * time.Millisecond

	timeout := time.NewTimer(duration)
	defer timeout.Stop()

	select {
	case <-w.ctx.Done():
		return false
	case <-timeout.C:
		return true
	}
}