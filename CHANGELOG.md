# Changelog

## [1.4.0](https://github.com/meysam81/parse-dmarc/compare/v1.3.10...v1.4.0) (2025-12-17)


### Features

* **docs:** add 1-click deployment at the top of README ([#41](https://github.com/meysam81/parse-dmarc/issues/41)) ([50d2ff2](https://github.com/meysam81/parse-dmarc/commit/50d2ff2c54a506801e802daa6b5d375b79d1281d))
* **docs:** add northflank deployment button ([81215a5](https://github.com/meysam81/parse-dmarc/commit/81215a5b5e8cd87aa5f601f22fb460b4929d8d15))
* **docs:** add the self-hosted options to deployments ([42ab966](https://github.com/meysam81/parse-dmarc/commit/42ab966218b9e875ed438a9cb4a4b4ac3adb9b33))
* **docs:** add zeabur deployment template ([4abaa36](https://github.com/meysam81/parse-dmarc/commit/4abaa36d9b077f71a3e6da841af056dffbe046ba))
* **UI:** revamp the dashboard for actionable insights ([#68](https://github.com/meysam81/parse-dmarc/issues/68)) ([dbe073d](https://github.com/meysam81/parse-dmarc/commit/dbe073d35a7f349770e9fab6df4923e221c350f2))


### Bug Fixes

* **deps:** update module github.com/coreos/go-oidc/v3 to v3.17.0 ([#64](https://github.com/meysam81/parse-dmarc/issues/64)) ([306583c](https://github.com/meysam81/parse-dmarc/commit/306583cc37ba2c4fc7e7fffd542cdbab3edd119d))
* **docs:** clean up the providers ([6879ca7](https://github.com/meysam81/parse-dmarc/commit/6879ca70d2a9e5d1adc7f7732d2a85869aecdd2e))
* **docs:** reference dockerhub image ([b2d2b06](https://github.com/meysam81/parse-dmarc/commit/b2d2b0627dbf423d760b5ee9598060334b13f290))
* update docker image reference to dockerhub ([6afcb3a](https://github.com/meysam81/parse-dmarc/commit/6afcb3abe19beda674d19ab7911019c384ba54b1))

## [1.3.10](https://github.com/meysam81/parse-dmarc/compare/v1.3.9...v1.3.10) (2025-12-13)


### Bug Fixes

* **CI:** reverse the digest conditional for provenance ([9c0e165](https://github.com/meysam81/parse-dmarc/commit/9c0e165ec1cfeb58c7c535ec02bd0c81e60ba501))

## [1.3.9](https://github.com/meysam81/parse-dmarc/compare/v1.3.8...v1.3.9) (2025-12-13)


### Bug Fixes

* **CI:** perform keyless docker sign ([1f322d3](https://github.com/meysam81/parse-dmarc/commit/1f322d33999d3586880d82c5ce3570b60157cd21))
* **CI:** update the digest name for attestation ([eb692a8](https://github.com/meysam81/parse-dmarc/commit/eb692a81fb08205990e639ab6c54c8998c705cd0))

## [1.3.8](https://github.com/meysam81/parse-dmarc/compare/v1.3.7...v1.3.8) (2025-12-13)


### Bug Fixes

* **CI:** set up buildx action for multi platform build ([a36e41c](https://github.com/meysam81/parse-dmarc/commit/a36e41c02e6af54297275042e269c66b85545a78))

## [1.3.7](https://github.com/meysam81/parse-dmarc/compare/v1.3.6...v1.3.7) (2025-12-13)


### Bug Fixes

* **CI:** remove annotations from docker build ([b5d0cd6](https://github.com/meysam81/parse-dmarc/commit/b5d0cd650ba0dd780ea95ec35c82940620aad6bb))

## [1.3.6](https://github.com/meysam81/parse-dmarc/compare/v1.3.5...v1.3.6) (2025-12-13)


### Bug Fixes

* **CI:** disable sbom via config ([6b0f3be](https://github.com/meysam81/parse-dmarc/commit/6b0f3be7d79cafe1eeb0930d2f30bf02f53c56af))

## [1.3.5](https://github.com/meysam81/parse-dmarc/compare/v1.3.4...v1.3.5) (2025-12-13)


### Bug Fixes

* **CI:** disable attestation digest for docker build ([95d6665](https://github.com/meysam81/parse-dmarc/commit/95d666543f59b3a890ea8921275f840eebb32f8e))

## [1.3.4](https://github.com/meysam81/parse-dmarc/compare/v1.3.3...v1.3.4) (2025-12-13)


### Bug Fixes

* **CI:** disable sbom on docker ([36e5cc4](https://github.com/meysam81/parse-dmarc/commit/36e5cc4d829efa633fc75879e0a132d076866318))

## [1.3.3](https://github.com/meysam81/parse-dmarc/compare/v1.3.2...v1.3.3) (2025-12-13)


### Bug Fixes

* **CI:** add current dir to docker build context ([254cc6a](https://github.com/meysam81/parse-dmarc/commit/254cc6a6b0b88afc67af8a1cdc4a41d1b71ff1a6))

## [1.3.2](https://github.com/meysam81/parse-dmarc/compare/v1.3.1...v1.3.2) (2025-12-13)


### Features

* **CI:** build docker via goreleaser ([8555265](https://github.com/meysam81/parse-dmarc/commit/855526505aa48c9ab1cdc946f10370a875ee47d8))

## [1.3.1](https://github.com/meysam81/parse-dmarc/compare/v1.3.0...v1.3.1) (2025-12-13)


### Bug Fixes

* **CI:** publish to casks for a change ([0e8346b](https://github.com/meysam81/parse-dmarc/commit/0e8346b3e1c3dcc557445370f59ea7c93588c1fd))

## [1.3.0](https://github.com/meysam81/parse-dmarc/compare/v1.2.2...v1.3.0) (2025-12-12)


### Features

* add MCP server integration to project ([#46](https://github.com/meysam81/parse-dmarc/issues/46)) ([807b8d6](https://github.com/meysam81/parse-dmarc/commit/807b8d677c81aef6ac39605bffb7125680841d4a))
* **CI:** add man page and shell completion to brew ([0e45ed2](https://github.com/meysam81/parse-dmarc/commit/0e45ed2a10d26a06fadc6a2fdfe3b81ad5a01197))
* **CI:** add prettier job ([39af52d](https://github.com/meysam81/parse-dmarc/commit/39af52d186fcf75dba2f8ac05ef313808f1e01b4))
* **CI:** install pandoc and use official goreleaser action ([0d0a307](https://github.com/meysam81/parse-dmarc/commit/0d0a3077395a9879a9d036166b4afbab43ee164e))
* **mcp:** add OAuth2 authentication for MCP HTTP server ([361f078](https://github.com/meysam81/parse-dmarc/commit/361f0782e2f20b7160873a085e600adc68988432))

## [1.2.2](https://github.com/meysam81/parse-dmarc/compare/v1.2.1...v1.2.2) (2025-12-05)


### Bug Fixes

* **CI:** add the brew token ([4081621](https://github.com/meysam81/parse-dmarc/commit/4081621ae9d1771185712512a912aa350a2e0305))

## [1.2.1](https://github.com/meysam81/parse-dmarc/compare/v1.2.0...v1.2.1) (2025-12-05)


### Bug Fixes

* **CI:** move bundled dist to server directory ([59dabf1](https://github.com/meysam81/parse-dmarc/commit/59dabf1becfc22fddda4bfdab6dd22188396b4c3))
* **docs:** pin to full version for now ([2060e45](https://github.com/meysam81/parse-dmarc/commit/2060e451a3b7322d933ad4b8f389f28474f85fb1))
* **docs:** update quickstart command with volume ([617c258](https://github.com/meysam81/parse-dmarc/commit/617c2584402cb309e1a652630f5f75de92df0b2e))

## [1.2.0](https://github.com/meysam81/parse-dmarc/compare/v1.1.0...v1.2.0) (2025-12-05)


### Features

* add production-ready Grafana dashboard ([6f857fe](https://github.com/meysam81/parse-dmarc/commit/6f857fe660144093b8766f529c15b17a21edc3b2))
* add production-ready Prometheus metrics ([#39](https://github.com/meysam81/parse-dmarc/issues/39)) ([17f6968](https://github.com/meysam81/parse-dmarc/commit/17f69681b23eee92ff1bc908fde7da2de11b8a52))
* **frontend:** add DMARC DNS record generator ([34491ba](https://github.com/meysam81/parse-dmarc/commit/34491ba244f9af48a25e3a0d6b09979661678c0d))
* **frontend:** implement dark mode with theme toggle ([#35](https://github.com/meysam81/parse-dmarc/issues/35)) ([569fd36](https://github.com/meysam81/parse-dmarc/commit/569fd3649931f94fa4f39845d0f141afacd6fbc0))


### Bug Fixes

* **CI:** build on major version as well ([fdff51d](https://github.com/meysam81/parse-dmarc/commit/fdff51d2c77d55e890e5b6d9c12ce316ddec2542))
* **CI:** only pin main with latest ([f120b2e](https://github.com/meysam81/parse-dmarc/commit/f120b2e7b7985adbaa815b6104f8aaee44e8e9ce))
* **CI:** reverse the conditional for build-dev job ([ce8b417](https://github.com/meysam81/parse-dmarc/commit/ce8b417c65c62563886787df3871200679c8de66))
* **CI:** update dockerignore after main.go change ([1241b2f](https://github.com/meysam81/parse-dmarc/commit/1241b2f22b8138348b7bf087ba00659e0b1f6fe8))
* **CI:** update version of the binary ([0a8f023](https://github.com/meysam81/parse-dmarc/commit/0a8f0230b30229349db270aa065e3332da63990a))
* **CI:** use go version file when setting up go ([40fbcb1](https://github.com/meysam81/parse-dmarc/commit/40fbcb1dfbe962c7d48ef2c821ee32d997e7c78c))
* **dev:** remove extra compose file ([eb71b09](https://github.com/meysam81/parse-dmarc/commit/eb71b097005bb9d024454ee6d84c2ce5fbb59d7b))
* **dev:** remove user from compose ([e9fbf73](https://github.com/meysam81/parse-dmarc/commit/e9fbf73b72a2a9958d62ed6f6ffea4bad2f2771f))
* make github star non-intrusive ([831047b](https://github.com/meysam81/parse-dmarc/commit/831047b88394cd468ff57bc08c53907880e0feee))
* move main.go to the root ([0dd5787](https://github.com/meysam81/parse-dmarc/commit/0dd57878ec5f7d6707faca0189834a6d91c266e5))
* publish to grafana community dashboards instead ([86f40b0](https://github.com/meysam81/parse-dmarc/commit/86f40b052607861de3d1a843f9e9c9e08ba57481))
* remove broken dns generator ([df39352](https://github.com/meysam81/parse-dmarc/commit/df393527710b487bc0bdbe7d67d9a3f909623b92))

## [1.1.0](https://github.com/meysam81/parse-dmarc/compare/v1.0.2...v1.1.0) (2025-11-08)


### Features

* **docs:** add nerdy badges to README ([#23](https://github.com/meysam81/parse-dmarc/issues/23)) ([5026a45](https://github.com/meysam81/parse-dmarc/commit/5026a45da22a9fa0a3c831945da4afaa05c2dd7c))


### Bug Fixes

* **CI:** update goreleaser after moving FE to root ([357d37e](https://github.com/meysam81/parse-dmarc/commit/357d37ed69d9cd7400c05c2f7ca2667072f4d46d))

## [1.0.2](https://github.com/meysam81/parse-dmarc/compare/v1.0.1...v1.0.2) (2025-11-07)


### Bug Fixes

* **dev:** bring the frontend to the root and simplify docker ([edc64e7](https://github.com/meysam81/parse-dmarc/commit/edc64e7417c910b7200bb1d727f0ba200c1a787d))

## [1.0.1](https://github.com/meysam81/parse-dmarc/compare/v1.0.0...v1.0.1) (2025-11-07)


### Bug Fixes

* **docs:** simplify README for first-time viewers ([#15](https://github.com/meysam81/parse-dmarc/issues/15)) ([7341b7b](https://github.com/meysam81/parse-dmarc/commit/7341b7bb022e5e96798e6d8d3523a377e282cf6f))

## 1.0.0 (2025-11-05)


### Features

* add CI and goreleaser with optimized distroless docker ([66c6828](https://github.com/meysam81/parse-dmarc/commit/66c682807d5ed12328349e2713c9d814d6337e78))
* add sort and refresh to UI ([9d6853e](https://github.com/meysam81/parse-dmarc/commit/9d6853e95319529c40958ae60d95d05d8cb5a675))
* clean up the cli and optimize the docker ([db23ba4](https://github.com/meysam81/parse-dmarc/commit/db23ba4ff2fd80ad1c016daa73a57519998dd271))
* **docs:** update readme to the latest changes ([c20e0cf](https://github.com/meysam81/parse-dmarc/commit/c20e0cf837b91b05e61153e4e55e7315fb09d396))
* ensure db path exists and add demo screenshot ([42d6d14](https://github.com/meysam81/parse-dmarc/commit/42d6d14fd96f2818830c73ad4e23b8f97719c497))
* Implement DMARC report parser with Vue.js dashboard ([#1](https://github.com/meysam81/parse-dmarc/issues/1)) ([23b8ac0](https://github.com/meysam81/parse-dmarc/commit/23b8ac0cc63a81aad7dfa2385ed24f6d74915775))
* update the UI footer with OS friendly text ([5008fdf](https://github.com/meysam81/parse-dmarc/commit/5008fdf5e59e5b4af1fc12ccafc644d27ccdbdfc))
* update UI footer & readme and optimize vite ([9578d05](https://github.com/meysam81/parse-dmarc/commit/9578d0530adab10334ad1a8a20c563abdc589854))


### Bug Fixes

* **CI:** ensure a fake dist exists before golangci-lint ([4e01983](https://github.com/meysam81/parse-dmarc/commit/4e01983f13a16ab0a8ea7cb880e43b39c285ba1f))
* **CI:** make linter happy ([cd04dfc](https://github.com/meysam81/parse-dmarc/commit/cd04dfc8beae30b0405d38834680011cd7d14fe5))
* make server context aware and do not fail UI on empty result ([588c466](https://github.com/meysam81/parse-dmarc/commit/588c4660534d4111cf24738a952cbad56ebbf345))
* update the UI on reports API parsing ([b85b59e](https://github.com/meysam81/parse-dmarc/commit/b85b59ef577ffeb6b4dacd1188e6e55d549b1b5a))
* use app for static config ([3820862](https://github.com/meysam81/parse-dmarc/commit/3820862f8fffe26fd3dad7c004d1c0b04fc8be9e))
