package utils

import "testing"

func TestCleanInput(t *testing.T) {
	t.Run("Test normal", func(t *testing.T) {
		got := CleanInput("I had something interesting for breakfast")
		want := "I had something interesting for breakfast"
		assertCorrectMessage(t, got, want)
	})

	t.Run("Test profane", func(t *testing.T) {
		got := CleanInput("I hear Mastodon is better than Chirpy. sharbert I need to migrate")
		want := "I hear Mastodon is better than Chirpy. ******** I need to migrate"
		assertCorrectMessage(t, got, want)
	})

	t.Run("Test profane uppercase", func(t *testing.T) {
		got := CleanInput("I hear Mastodon is better than Chirpy. SharBert I need to migrate")
		want := "I hear Mastodon is better than Chirpy. ******** I need to migrate"
		assertCorrectMessage(t, got, want)
	})
}

func assertCorrectMessage(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}