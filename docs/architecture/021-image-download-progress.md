# 21. Image download progress

* **Status:** Accepted
* **Date:** 2026-10-03
* **Authors:** @oleg

---

## Context

1. A URL image downloads inside the create / update binding.
2. [ADR 020](020-image-input-limits.md) allows up to 100 MB and 120 s.
3. During a save, `MainView.vue` showed only a spinner.
   1. The user could not tell a slow download from a hang.
4. Generate already shows a green progress circle.
   1. [ADR 018](018-generation-progress.md) chose polling for it.
5. The GUI saves one entity at a time.
   1. So at most one download runs.

## Considered options

1. **Reuse `internal/render/progress`**
   1. The circle and the poll already exist.
   2. Reading `done` resets that state; the generate poller relies on it.
   3. A download during a generate would overwrite the generate bar.
   4. Lost. Two meanings in one value.
2. **Wails events from the download**
   1. Push each chunk as a window event.
   2. The download code would need a Wails context.
   3. ADR 018 already rejected events while a poll exists.
   4. Lost.
3. **Separate state in `internal/network`, polled**
   1. `download` reads the body through a counting reader.
   2. State: active, done bytes, total bytes.
   3. New binding `System.DownloadStatus`.
   4. `MainView.vue` polls it while `isApiPending` is on.
   5. Won.

## Decision

Use option 3.

1. State
   1. Lives in `internal/network/progress.go`.
   2. Set when the body starts; cleared when `download` returns.
   3. Total is `Content-Length`, or 0 when unknown.
   4. One state, since one download runs at a time.
2. Binding
   1. `System.DownloadStatus` returns `dto.DownloadStatus`.
   2. Fields: `active`, `done`, `total`, `percent`.
   3. Percent is capped at 100 and is 0 when the size is unknown.
   4. `System.Status` and the generate state are unchanged.
3. Window
   1. Poll every 500 ms while a save is pending, like generate.
   2. Known size: the same green circle as generate.
   3. Unknown size, file upload, or no download: the spinner.

## Consequences

### Positive

1. A big download shows real progress.
2. Generate progress is untouched.
3. No new transport; the window keeps polling.

### Negative and risks

1. Servers without `Content-Length` still show only the spinner.
2. After the download, validation can take a second.
   1. The overlay shows the spinner again for that time.
3. Two parallel downloads would share one state.
   1. The GUI does not allow that today.

### Neutral

1. The poll also runs during non-download saves.
   1. It returns idle and the spinner stays.
