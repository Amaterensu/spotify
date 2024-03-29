package spotify

import (
	"context"
	"net/http"
	"testing"
)

// The example from https://developer.spotify.com/web-api/get-album/
func TestFindAlbum(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/find_album.txt")
	defer server.Close()

	album, err := client.GetAlbum(context.Background(), ID("0sNOF9WDwhWunNAHPD3Baj"))
	if err != nil {
		t.Fatal(err)
	}
	if album == nil {
		t.Fatal("Got nil album")
	}
	if album.Name != "She's So Unusual" {
		t.Error("Got wrong album")
	}
	released := album.ReleaseDateTime()
	if released.Year() != 1983 {
		t.Errorf("Expected release date 1983, got %d\n", released.Year())
	}
}

func TestFindAlbumBadID(t *testing.T) {
	client, server := testClientString(http.StatusNotFound, `{ "error": { "status": 404, "message": "non existing id" } }`)
	defer server.Close()

	album, err := client.GetAlbum(context.Background(), ID("asdf"))
	if album != nil {
		t.Fatal("Expected nil album, got", album.Name)
	}
	se, ok := err.(Error)
	if !ok {
		t.Error("Expected spotify error, got", err)
	}
	if se.Status != 404 {
		t.Errorf("Expected HTTP 404, got %d. ", se.Status)
	}
	if se.Message != "non existing id" {
		t.Error("Unexpected error message: ", se.Message)
	}
}

// The example from https://developer.spotify.com/web-api/get-several-albums/
func TestFindAlbums(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/find_albums.txt")
	defer server.Close()

	res, err := client.GetAlbums(context.Background(), []ID{"41MnTivkwTO3UUJ8DrqEJJ", "6JWc4iAiJ9FjyK0B59ABb4", "6UXCm6bOO4gFlDQZV5yL37", "0X8vBD8h1Ga9eLT8jx9VCC"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 4 {
		t.Fatalf("Expected 4 albums, got %d", len(res))
	}
	expectedAlbums := []string{
		"The Best Of Keane (Deluxe Edition)",
		"Strangeland",
		"Night Train",
		"Mirrored",
	}
	for i, name := range expectedAlbums {
		if res[i].Name != name {
			t.Error("Expected album", name, "but got", res[i].Name)
		}
	}
	release := res[0].ReleaseDateTime()
	if release.Year() != 2013 ||
		release.Month() != 11 ||
		release.Day() != 8 {
		t.Errorf("Expected release 2013-11-08, got %d-%02d-%02d\n",
			release.Year(), release.Month(), release.Day())
	}
	releaseMonthPrecision := res[3].ReleaseDateTime()
	if releaseMonthPrecision.Year() != 2007 ||
		releaseMonthPrecision.Month() != 3 ||
		releaseMonthPrecision.Day() != 1 {
		t.Errorf("Expected release 2007-03-01, got %d-%02d-%02d\n",
			releaseMonthPrecision.Year(), releaseMonthPrecision.Month(), releaseMonthPrecision.Day())
	}
}

func TestFindAlbumTracks(t *testing.T) {
	client, server := testClientFile(http.StatusOK, "test_data/find_album_tracks.txt")
	defer server.Close()

	res, err := client.GetAlbumTracks(context.Background(), ID("0sNOF9WDwhWunNAHPD3Baj"), Limit(1))
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 13 {
		t.Fatal("Got", res.Total, "results, want 13")
	}
	if len(res.Tracks) == 1 {
		if res.Tracks[0].Name != "Money Changes Everything" {
			t.Error("Expected track 'Money Changes Everything', got", res.Tracks[0].Name)
		}
	} else {
		t.Error("Expected 1 track, got", len(res.Tracks))
	}
}
