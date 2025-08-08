package session

import (
	"reflect"
	"testing"
)

var (
	mockSession1 = Session{
		Name:        "test-session-1",
		CurrentPath: "/home/user/project1",
		Windows:     []Window{},
	}
	mockSession2 = Session{
		Name:        "test-session-2",
		CurrentPath: "/home/user/project2",
		Windows:     []Window{},
	}
	mockSession3 = Session{
		Name:        "test-session-3",
		CurrentPath: "/home/user/project1",
		Windows:     []Window{},
	}
)

func TestSessionFunctions(t *testing.T) {
	sessions := []Session{mockSession1, mockSession2}

	t.Run("GetSessionByName", func(t *testing.T) {
		s, err := GetSessionByName("test-session-1", sessions)
		if err != nil {
			t.Errorf("expected session 'test-session-1' to be found, got error: %v", err)
		}
		if s.Name != "test-session-1" {
			t.Errorf("expected session name 'test-session-1', got '%s'", s.Name)
		}

		_, err = GetSessionByName("non-existent-session", sessions)
		if err == nil {
			t.Errorf("expected an error for non-existent session, got nil")
		}
	})

	t.Run("GetSessionByPath", func(t *testing.T) {
		s, err := GetSessionByPath("/home/user/project2", sessions)
		if err != nil {
			t.Errorf("expected session at '/home/user/project2' to be found, got error: %v", err)
		}
		if s.Name != "test-session-2" {
			t.Errorf("expected session name 'test-session-2', got '%s'", s.Name)
		}

		_, err = GetSessionByPath("/home/user/non-existent-path", sessions)
		if err == nil {
			t.Errorf("expected an error for non-existent path, got nil")
		}
	})

	t.Run("CombineSessions", func(t *testing.T) {
		newSessions := []Session{mockSession2, mockSession3}
		combined, err := CombineSessions(sessions, newSessions)
		if err != nil {
			t.Fatalf("CombineSessions failed: %v", err)
		}

		expected := []Session{mockSession1, mockSession2, mockSession3}
		if !reflect.DeepEqual(combined, expected) {
			t.Errorf("CombineSessions did not produce the expected result.\nExpected: %v\nGot:      %v", expected, combined)
		}
	})

	t.Run("DeleteSession", func(t *testing.T) {
		updatedSessions, err := DeleteSession("test-session-1", sessions)
		if err != nil {
			t.Fatalf("DeleteSession failed: %v", err)
		}

		expected := []Session{mockSession2}
		if !reflect.DeepEqual(updatedSessions, expected) {
			t.Errorf("DeleteSession did not remove the correct session.\nExpected: %v\nGot:      %v", expected, updatedSessions)
		}

		_, err = DeleteSession("non-existent", sessions)
		if err == nil {
			t.Errorf("expected an error for deleting a non-existent session, got nil")
		}
	})

	t.Run("CheckIfSessionExists", func(t *testing.T) {
		if !CheckIfSessionExists("test-session-2", sessions) {
			t.Errorf("expected 'test-session-2' to exist, but it was not found")
		}

		if CheckIfSessionExists("non-existent", sessions) {
			t.Errorf("expected 'non-existent' to not exist, but it was found")
		}
	})
}
