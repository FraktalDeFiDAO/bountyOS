# Scanner Migration Plan

To maintain the architectural integrity of BountyOS, all scanners must be migrated to the `BaseScanner` pattern. This pattern provides standardized HTTP clients, rate limiting, and result normalization.

## Completed Migrations
- [x] GitHub
- [x] Superteam
- [x] Bountycaster
- [x] Immunefi
- [x] Algora (Migrated 2026-02-27)

## Remaining Migrations
- [ ] Apify (`apify.go`)
- [ ] Base (`base.go`)
- [ ] CharmVerse (`charmverse.go`)
- [ ] Clawlancer (`clawlancer.go`)
- [ ] Code4rena (`code4rena.go`)
- [ ] Colosseum (`colosseum.go`)
- [ ] Dework (`dework.go`)
- [ ] LaborX (`laborx.go`)
- [ ] Optimism (`optimism.go`)
- [ ] Proxies (`proxies.go`)
- [ ] ugig (`ugig.go`)
- [ ] Uniswap (`uniswap.go`)

## Migration Guide
1. Embed `BaseScanner` in the scanner struct.
2. Use `NewBaseScanner` in the constructor.
3. Replace manual `http.Get` calls with `s.FetchJSON`.
4. Replace manual `core.Bounty` construction with `s.CreateBounty`.
5. Use `FormatAmount()` for reward strings.
