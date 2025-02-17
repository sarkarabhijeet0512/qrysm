package synccommittee

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/theQRL/qrysm/v4/consensus-types/primitives"
	"github.com/theQRL/qrysm/v4/container/queue"
	zondpb "github.com/theQRL/qrysm/v4/proto/prysm/v1alpha1"
)

// To give two slots tolerance for objects that arrive earlier.
// This account for previous slot, current slot, two future slots.
const syncCommitteeMaxQueueSize = 4

// SaveSyncCommitteeContribution saves a sync committee contribution in to a priority queue.
// The priority queue is capped at syncCommitteeMaxQueueSize contributions.
func (s *Store) SaveSyncCommitteeContribution(cont *zondpb.SyncCommitteeContribution) error {
	if cont == nil {
		return errNilContribution
	}

	s.contributionLock.Lock()
	defer s.contributionLock.Unlock()

	item, err := s.contributionCache.PopByKey(syncCommitteeKey(cont.Slot))
	if err != nil {
		return err
	}

	copied := zondpb.CopySyncCommitteeContribution(cont)

	// Contributions exist in the queue. Append instead of insert new.
	if item != nil {
		contributions, ok := item.Value.([]*zondpb.SyncCommitteeContribution)
		if !ok {
			return errors.New("not typed []zondpb.SyncCommitteeContribution")
		}

		contributions = append(contributions, copied)
		savedSyncCommitteeContributionTotal.Inc()
		return s.contributionCache.Push(&queue.Item{
			Key:      syncCommitteeKey(cont.Slot),
			Value:    contributions,
			Priority: int64(cont.Slot),
		})
	}

	// Contribution does not exist. Insert new.
	if err := s.contributionCache.Push(&queue.Item{
		Key:      syncCommitteeKey(cont.Slot),
		Value:    []*zondpb.SyncCommitteeContribution{copied},
		Priority: int64(cont.Slot),
	}); err != nil {
		return err
	}
	savedSyncCommitteeContributionTotal.Inc()

	// Trim contributions in queue down to syncCommitteeMaxQueueSize.
	if s.contributionCache.Len() > syncCommitteeMaxQueueSize {
		if _, err := s.contributionCache.Pop(); err != nil {
			return err
		}
	}

	return nil
}

// SyncCommitteeContributions returns sync committee contributions by slot from the priority queue.
// Upon retrieval, the contribution is removed from the queue.
func (s *Store) SyncCommitteeContributions(slot primitives.Slot) ([]*zondpb.SyncCommitteeContribution, error) {
	s.contributionLock.RLock()
	defer s.contributionLock.RUnlock()

	item := s.contributionCache.RetrieveByKey(syncCommitteeKey(slot))
	if item == nil {
		return []*zondpb.SyncCommitteeContribution{}, nil
	}

	contributions, ok := item.Value.([]*zondpb.SyncCommitteeContribution)
	if !ok {
		return nil, errors.New("not typed []zondpb.SyncCommitteeContribution")
	}

	return contributions, nil
}

func syncCommitteeKey(slot primitives.Slot) string {
	return strconv.FormatUint(uint64(slot), 10)
}
