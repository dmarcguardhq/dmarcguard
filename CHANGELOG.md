# Changelog

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
