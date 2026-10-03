# UC-13: Show image download progress

**Module:** `internal/network`, `bindings/system`  
**Status:** Implemented  
**Actors:** User saving an entity with an image URL; `MainView.vue` polling `System.DownloadStatus`  
**Goal:** While a big image downloads, the save overlay shows how far it got  
**Preconditions:** A create or update with a changed image URL (UC-01–UC-04)

## Main scenario (happy path)

1. The user saves an entity with a new image URL.
2. `itemsStore.isApiPending` turns on; `MainView.vue` polls `System.DownloadStatus` every 500 ms.
3. `network.DownloadBytes` starts the download and records `Content-Length` as the total.
4. Each read adds to `done`; the binding returns `percent = done / total × 100`.
5. The overlay shows the same green circle as generate, with that percent.
6. The download ends; the state goes back to idle. Validation runs; the overlay shows the spinner again.
7. The save returns; `isApiPending` turns off and polling stops.

## Alternative scenarios and errors

* **3a. No `Content-Length`:** `total` is 0; the overlay keeps the spinner.
* **3b. File upload, or an unchanged URL:** nothing downloads; the overlay keeps the spinner.
* **4a. Download fails (timeout, size, status):** the state goes back to idle; the save still succeeds with a `warning` (UC-01–UC-04).

## Postconditions

* After every download, success or not, `System.DownloadStatus` reports idle.
* Only one download runs at a time: the GUI saves one entity at a time.
