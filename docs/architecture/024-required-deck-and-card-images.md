# 24. Decks and cards require an image

* **Status:** Accepted
* **Date:** 2026-10-04
* **Authors:** @oleg

---

## Context

1. Rendering needs a back image per deck and a face per card.
   1. TTS cannot show a sheet without them.
2. Game and collection images only decorate the app.
3. Before this decision a missing image surfaced only at render time.
   1. The render stopped in its goroutine and only logged the cause.
   2. The window hid the progress circle and showed no reason.
4. Saving a deck or card without an image gave no hint at all.

## Considered options

1. **Optional images, with a "remove image" control**
   1. Asked for in review item S12.
   2. Lost. A deck or card without an image cannot be rendered.
2. **Block the save until an image is set**
   1. Lost. Users fill catalogs in steps; a hard stop gets in the way.
3. **Required, reminded softly, checked before rendering**
   1. Won.

## Decision

1. Scope
   1. Decks and cards must have an image.
   2. Games and collections may not; they get no reminders.
2. Saving
   1. A save without an image succeeds.
   2. The binding result carries a warning, shown as a toast.
   3. Deck: "This deck has no image. Rendering needs a back image for every deck."
   4. Card: "This card has no image. Rendering needs an image for every card."
3. Lists
   1. Deck and card DTOs carry `hasImage`.
   2. It comes from a folder listing; no image is read.
   3. The tile shows a "No image" placeholder and a warning badge.
   4. A deck also gets the badge when any of its cards has no image.
      1. `cardsMissingImage` in the deck DTO; the cards are not named.
      2. No roll-up to collections or games.
4. Rendering
   1. `GenerateGame` checks every deck back and card face first (`generate.MissingImages`).
   2. Anything missing: `GenerateMissingImages`, naming up to 5 and counting the rest.
   3. Nothing is written; the render does not start.
   4. A deck's back is the one in its first card's collection, as `Prepare` uses.
5. There is no control to remove an image.

## Consequences

### Positive

1. Users see what blocks a render before starting it.
2. A render never ends silently because of a missing image.
3. Partial TTS results with missing cards cannot happen.

### Negative and risks

1. Each deck and card list does one extra folder listing.
2. An image file that exists but cannot be decoded still fails during the render.
   1. That failure is still not shown in the window (review item D).

### Neutral

1. Review item S12 (remove a file-uploaded image) is closed by this decision.
