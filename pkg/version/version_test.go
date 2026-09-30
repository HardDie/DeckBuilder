package version

import (
	"os/exec"
	"strings"
	"testing"
)

func TestString_blankIsDev(t *testing.T) {
	old := Build
	t.Cleanup(func() { Build = old })

	Build = " "
	if String() != "dev" {
		t.Fatalf("got %q", String())
	}
	Build = "v1.2.3"
	if String() != "v1.2.3" {
		t.Fatalf("got %q", String())
	}
}

func TestString_defaultDev(t *testing.T) {
	if String() != "dev" {
		t.Fatalf("version = %q", String())
	}
}

func TestDescribe_tagOrHash(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	git(t, dir, "init", "-b", "main")
	git(t, dir, "config", "user.email", "test@example.com")
	git(t, dir, "config", "user.name", "Test")
	git(t, dir, "config", "commit.gpgsign", "false")
	git(t, dir, "commit", "--allow-empty", "-m", "init")

	got, err := Describe(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := gitOut(t, dir, "rev-parse", "--short=12", "HEAD")
	if got != want {
		t.Fatalf("hash: got %q want %q", got, want)
	}
	if len(got) != 12 {
		t.Fatalf("hash length %d", len(got))
	}

	git(t, dir, "tag", "v9.9.9")
	got, err = Describe(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "v9.9.9" {
		t.Fatalf("tag: got %q", got)
	}

	git(t, dir, "commit", "--allow-empty", "-m", "second")
	got, err = Describe(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got == "v9.9.9" {
		t.Fatal("untagged commit kept the old tag")
	}
	want = gitOut(t, dir, "rev-parse", "--short=12", "HEAD")
	if got != want {
		t.Fatalf("later hash: got %q want %q", got, want)
	}
}

func TestDescribe_notRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	if _, err := Describe(t.TempDir()); err == nil {
		t.Fatal("expected error")
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}

func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}
