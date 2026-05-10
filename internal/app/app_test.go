package app

import (
	"context"
	"testing"

	"backlogs-bot/internal/seatalk"
)

type fakeGroupIDStore struct {
	upsertTab string
	upsertID  string
	removeTab string
	removeID  string
}

func (f *fakeGroupIDStore) UpsertGroupID(_ context.Context, tab, groupID string) error {
	f.upsertTab = tab
	f.upsertID = groupID
	return nil
}

func (f *fakeGroupIDStore) RemoveGroupID(_ context.Context, tab, groupID string) error {
	f.removeTab = tab
	f.removeID = groupID
	return nil
}

func TestHandleSeaTalkEventStoresJoinedGroupID(t *testing.T) {
	var event seatalk.CallbackEvent
	event.EventType = seatalk.EventBotAddedToGroupChat
	event.Event.Group.GroupID = "group-123"
	event.Event.Group.GroupName = "SOC Alerts"

	store := &fakeGroupIDStore{}
	if err := handleSeaTalkEvent(context.Background(), event, "BAU Backlogs Summary", store); err != nil {
		t.Fatalf("handle event: %v", err)
	}

	if store.upsertTab != "BAU Backlogs Summary" {
		t.Fatalf("upsert tab = %q, want BAU Backlogs Summary", store.upsertTab)
	}
	if store.upsertID != "group-123" {
		t.Fatalf("upsert group ID = %q, want group-123", store.upsertID)
	}
}

func TestHandleSeaTalkEventRemovesGroupID(t *testing.T) {
	var event seatalk.CallbackEvent
	event.EventType = seatalk.EventBotRemovedFromGroupChat
	event.Event.Group.GroupID = "group-123"

	store := &fakeGroupIDStore{}
	if err := handleSeaTalkEvent(context.Background(), event, "BAU Backlogs Summary", store); err != nil {
		t.Fatalf("handle event: %v", err)
	}

	if store.removeTab != "BAU Backlogs Summary" {
		t.Fatalf("remove tab = %q, want BAU Backlogs Summary", store.removeTab)
	}
	if store.removeID != "group-123" {
		t.Fatalf("remove group ID = %q, want group-123", store.removeID)
	}
}
