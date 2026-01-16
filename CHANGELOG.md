# Changelog

## [1.0.2-alpha.2](https://github.com/magenx/hek3ster/compare/v1.0.2-alpha.1...v1.0.2-alpha.2) (2026-01-16)


### 🚦 Maintenance

* Modify build workflow for draft releases ([#124](https://github.com/magenx/hek3ster/issues/124)) ([6992579](https://github.com/magenx/hek3ster/commit/699257982c3d80d35fe4b544a7675ddd7adf4fb2))
* **release:** Add draft option to prerelease configuration ([#122](https://github.com/magenx/hek3ster/issues/122)) ([ace11d3](https://github.com/magenx/hek3ster/commit/ace11d3624302d65b46cba16a492abc959adc9a7))
* Remove commented-out code in init function ([#125](https://github.com/magenx/hek3ster/issues/125)) ([4f781da](https://github.com/magenx/hek3ster/commit/4f781da51a9e40175be08ebea0f0fe2bb24951da))

## [1.0.2-alpha.1](https://github.com/magenx/hek3ster/compare/v1.0.2-alpha...v1.0.2-alpha.1) (2026-01-16)


### 🐛 Bug Fixes

* **addons:** Fix Cilium CNI installation failure with tilde-prefixed kubeconfig paths ([#117](https://github.com/magenx/hek3ster/issues/117)) ([9f15d14](https://github.com/magenx/hek3ster/commit/9f15d14327b8e9572c2cfcea5f1bf6dc88b9a939))
* **load balancer:** Fix load balancer network attachment validation for private IP targets ([#120](https://github.com/magenx/hek3ster/issues/120)) ([cf0108b](https://github.com/magenx/hek3ster/commit/cf0108be1732b2a8f1f59ca05c98f9fe5ff8bd73))

## [1.0.2-alpha](https://github.com/magenx/hek3ster/compare/v1.0.1...v1.0.2-alpha) (2026-01-16)


### 🐛 Bug Fixes

* **addons:** Change HubbleMetrics type from string to []string ([#109](https://github.com/magenx/hek3ster/issues/109)) ([80d2823](https://github.com/magenx/hek3ster/commit/80d2823c17869b51834a85759e94c413796bb915))
* **addons:** Replace Helm-based Cilium installation with official Cilium CLI approach ([#108](https://github.com/magenx/hek3ster/issues/108)) ([e5eb7fe](https://github.com/magenx/hek3ster/commit/e5eb7fe258bcfa681a8a98aa5e531e3eff69d492))


### 🚦 Maintenance

* **copilot:** Add GitHub Copilot custom instructions for repository context ([#111](https://github.com/magenx/hek3ster/issues/111)) ([bc2d1fd](https://github.com/magenx/hek3ster/commit/bc2d1fd2c1ff291f1b015ec16bc705e650d7ce58))
* **deps:** bump @radix-ui/react-toast from 1.2.14 to 1.2.15 in /page ([#88](https://github.com/magenx/hek3ster/issues/88)) ([824faea](https://github.com/magenx/hek3ster/commit/824faea292e690a571e30ba4ffe848144ec43a37))
* **deps:** bump next-themes from 0.3.0 to 0.4.6 in /page ([#87](https://github.com/magenx/hek3ster/issues/87)) ([1b11ef2](https://github.com/magenx/hek3ster/commit/1b11ef2463ae49025b46b782dd91f9f8f77852e9))
* **page:** Add HardDrive and ShieldCheck icons to Architecture ([#77](https://github.com/magenx/hek3ster/issues/77)) ([5a16d47](https://github.com/magenx/hek3ster/commit/5a16d4725545d58a0aabd58324d3077f8123d78c))
* **page:** Add links to github wiki ([#76](https://github.com/magenx/hek3ster/issues/76)) ([ffe7f3a](https://github.com/magenx/hek3ster/commit/ffe7f3ac6c1aa7de8ddfeec1010f5a1f30b059cb))
* **page:** Add linting and dependency check to deploy workflow ([#91](https://github.com/magenx/hek3ster/issues/91)) ([5a79a8b](https://github.com/magenx/hek3ster/commit/5a79a8bda4119f1be23727b01f0e0e0e6375eeaa))
* **page:** Add npm updates to dependabot configuration ([#85](https://github.com/magenx/hek3ster/issues/85)) ([efab961](https://github.com/magenx/hek3ster/commit/efab9615ec9d772225faee42f0f06486b628b6c7))
* **page:** Allow Check Missing Imports to continue on error ([#102](https://github.com/magenx/hek3ster/issues/102)) ([e32e65f](https://github.com/magenx/hek3ster/commit/e32e65fe4274e6f201372c7ea48660f6f46bf093))
* **page:** Change issue assignee from copilot to me ([#98](https://github.com/magenx/hek3ster/issues/98)) ([0f00725](https://github.com/magenx/hek3ster/commit/0f007252cb619fb558aad5eec2f4dcb2e8576f8c))
* **page:** Enable cancellation of in-progress page deployments ([#104](https://github.com/magenx/hek3ster/issues/104)) ([e154567](https://github.com/magenx/hek3ster/commit/e15456777c5f1f19999bc263429a2000c6c6865c))
* **page:** Enhance CI with linting and depcheck error handling ([#93](https://github.com/magenx/hek3ster/issues/93)) ([4b2d09d](https://github.com/magenx/hek3ster/commit/4b2d09dc4d0a199f901911fbbf0e2a122ba40942))
* **page:** Enhance Global Load Balancer with Firewall label ([#92](https://github.com/magenx/hek3ster/issues/92)) ([145df48](https://github.com/magenx/hek3ster/commit/145df48bba5fc6fdeb58915c20050a55818ebf30))
* **page:** Fix lint errors in UI components ([#97](https://github.com/magenx/hek3ster/issues/97)) ([a2282e1](https://github.com/magenx/hek3ster/commit/a2282e10fe6ab1c8a9cbcc7152fce39faa520eef))
* **page:** Modify hek3ster installation code snippet ([#75](https://github.com/magenx/hek3ster/issues/75)) ([e40d6e7](https://github.com/magenx/hek3ster/commit/e40d6e7e78a725c6cb554d072a519348c750e288))
* **page:** Remove unused dependencies from page/ ([#100](https://github.com/magenx/hek3ster/issues/100)) ([a74142d](https://github.com/magenx/hek3ster/commit/a74142d830c7ebeacc97a9784b7751853d5bdf80))
* **page:** Replace Github icon import with ExternalLink ([#83](https://github.com/magenx/hek3ster/issues/83)) ([0134aed](https://github.com/magenx/hek3ster/commit/0134aed8baad0ca93632da7c7dc2f2d60b0adfdd))
* **page:** Update documentation link in Hero component ([#84](https://github.com/magenx/hek3ster/issues/84)) ([9aaba59](https://github.com/magenx/hek3ster/commit/9aaba594c171fffa345af1f63d532c6dce6a13f2))
* **page:** Update footer link to new documentation site ([#82](https://github.com/magenx/hek3ster/issues/82)) ([4c6dfdb](https://github.com/magenx/hek3ster/commit/4c6dfdbffb7d6c7cb08a17811b37a10b565931ba))
* **page:** Update log formatting in deploy-pages.yml ([#95](https://github.com/magenx/hek3ster/issues/95)) ([3fdd967](https://github.com/magenx/hek3ster/commit/3fdd967bc88c77478904df5601140b0f17800b72))
* **page:** Update pricing details and improve descriptions ([#79](https://github.com/magenx/hek3ster/issues/79)) ([ae1e7b9](https://github.com/magenx/hek3ster/commit/ae1e7b928802dab9bb83b3e381121efa86948e48))
* **page:** Update pricing per custom setup upto date ([#78](https://github.com/magenx/hek3ster/issues/78)) ([8cb4cde](https://github.com/magenx/hek3ster/commit/8cb4cde0998982201c0f8a8e517d1ae5454d1148))
* **page:** Update savings percentage from 80% to 85% ([#80](https://github.com/magenx/hek3ster/issues/80)) ([a343a2a](https://github.com/magenx/hek3ster/commit/a343a2a90de37b010aea6ed4136eeab861cdf258))
* **workflow:** Rename build-pre-release.yml to build_for_release.yml ([#81](https://github.com/magenx/hek3ster/issues/81)) ([930d357](https://github.com/magenx/hek3ster/commit/930d357819e2f7ee103455e0bff3619eff5789dc))

## [1.0.1](https://github.com/magenx/hek3ster/compare/v1.0.0...v1.0.1) (2026-01-14)


### 🚦 Maintenance

* Add comment placeholder in init function ([#34](https://github.com/magenx/hek3ster/issues/34)) ([98803ca](https://github.com/magenx/hek3ster/commit/98803ca0113650164441b39c31308834e9a4ac10))
* Add comment placeholder in init function ([#43](https://github.com/magenx/hek3ster/issues/43)) ([a0e481a](https://github.com/magenx/hek3ster/commit/a0e481ac4a5febf30cb0a05ec2dcaab16edaf1ea))
* Add comment placeholder in init function ([#51](https://github.com/magenx/hek3ster/issues/51)) ([604c36a](https://github.com/magenx/hek3ster/commit/604c36a921e549839294cd50714734f847380667))
* Add comment to init function in completion.go ([#58](https://github.com/magenx/hek3ster/issues/58)) ([742a763](https://github.com/magenx/hek3ster/commit/742a76369e0b7638ded7d35436d20c097841707b))
* Add GH_TOKEN to build environment ([#18](https://github.com/magenx/hek3ster/issues/18)) ([d4a31d9](https://github.com/magenx/hek3ster/commit/d4a31d9c09acd5fd63227a945697b56e1ffd9737))
* Add GH_TOKEN to build-release workflow ([#19](https://github.com/magenx/hek3ster/issues/19)) ([a6547c9](https://github.com/magenx/hek3ster/commit/a6547c9a2522b65fd083c3c40605ecf5982dd442))
* Add permissions for contents in build workflow ([#32](https://github.com/magenx/hek3ster/issues/32)) ([c25f1c2](https://github.com/magenx/hek3ster/commit/c25f1c273b635d99e496c2d5f07f45abfd1148a7))
* Build workflow update ([da43e6c](https://github.com/magenx/hek3ster/commit/da43e6c6de9d7cdd0caff0e4910fbaa89097918c))
* Build workflow update release pre release ([2c1cd97](https://github.com/magenx/hek3ster/commit/2c1cd975d386fef52ceaa0b42fa85ad95766b412))
* Change cancel-in-progress to false in workflow ([#24](https://github.com/magenx/hek3ster/issues/24)) ([fb7dc90](https://github.com/magenx/hek3ster/commit/fb7dc90fd762403aa346af410c49e79161e42ae1))
* Delete .github/workflows/build-release.yml ([#66](https://github.com/magenx/hek3ster/issues/66)) ([4ef0868](https://github.com/magenx/hek3ster/commit/4ef0868b33782b443e0cc95e107281977f99c375))
* Delete CODE_OF_CONDUCT.md ([#20](https://github.com/magenx/hek3ster/issues/20)) ([36d35c5](https://github.com/magenx/hek3ster/commit/36d35c5e4a8e9d087f743a2ef19ad7f4ef503e9d))
* Enhance build script with binary checks and uploads ([#31](https://github.com/magenx/hek3ster/issues/31)) ([c014b2b](https://github.com/magenx/hek3ster/commit/c014b2ba84ed3971abf55299e18144a37d816249))
* Enhance build-pre-release.yml for artifact uploads ([#42](https://github.com/magenx/hek3ster/issues/42)) ([0062aa1](https://github.com/magenx/hek3ster/commit/0062aa132ef155bad5d8238f44b833464eeb6bff))
* Enhance installation summary formatting in YAML ([#65](https://github.com/magenx/hek3ster/issues/65)) ([35f322c](https://github.com/magenx/hek3ster/commit/35f322c240a32c7cb0a56e3031baca4498b4b13d))
* Exclude some files from workflow  trigger ([2bd0504](https://github.com/magenx/hek3ster/commit/2bd050422a75b4f6af32e8ed910358e7fa16e15c))
* Fix build-all target in Makefile ([#25](https://github.com/magenx/hek3ster/issues/25)) ([d6a7ce4](https://github.com/magenx/hek3ster/commit/d6a7ce4802cc4fab5132499b179d503b452a0592))
* Fix formatting issues in completion.go ([#71](https://github.com/magenx/hek3ster/issues/71)) ([4a25c47](https://github.com/magenx/hek3ster/commit/4a25c47afa43400452626de745b850ed6154378f))
* **main:** hek3ster 1.0.1-alpha ([#12](https://github.com/magenx/hek3ster/issues/12)) ([13b137d](https://github.com/magenx/hek3ster/commit/13b137d2eb4332d6d7b1881996bc4525615a8c39))
* **main:** hek3ster 1.0.1-alpha.1 ([#15](https://github.com/magenx/hek3ster/issues/15)) ([926b237](https://github.com/magenx/hek3ster/commit/926b237eac40bfc57ed62f7185d6990237ac5583))
* **main:** hek3ster 1.0.1-alpha.10 ([#55](https://github.com/magenx/hek3ster/issues/55)) ([8024abc](https://github.com/magenx/hek3ster/commit/8024abc421ecc60ee42017344021eed932a8c108))
* **main:** hek3ster 1.0.1-alpha.11 ([#59](https://github.com/magenx/hek3ster/issues/59)) ([5077ad7](https://github.com/magenx/hek3ster/commit/5077ad78cac7583d4759af019bae1e469c413028))
* **main:** hek3ster 1.0.1-alpha.12 ([#63](https://github.com/magenx/hek3ster/issues/63)) ([59775ed](https://github.com/magenx/hek3ster/commit/59775ed6bf6459be097501cd56a4c8684ffcaed5))
* **main:** hek3ster 1.0.1-alpha.13 ([#72](https://github.com/magenx/hek3ster/issues/72)) ([0b0cf90](https://github.com/magenx/hek3ster/commit/0b0cf9090d638cd3e47bac199d25783c60732d22))
* **main:** hek3ster 1.0.1-alpha.2 ([#22](https://github.com/magenx/hek3ster/issues/22)) ([c319940](https://github.com/magenx/hek3ster/commit/c31994006abcfef4508695547127f2fb3a9927a2))
* **main:** hek3ster 1.0.1-alpha.3 ([#26](https://github.com/magenx/hek3ster/issues/26)) ([7afba59](https://github.com/magenx/hek3ster/commit/7afba5910ecbb16e860b1a2ad4beee7e17ea4519))
* **main:** hek3ster 1.0.1-alpha.4 ([#29](https://github.com/magenx/hek3ster/issues/29)) ([daf8a65](https://github.com/magenx/hek3ster/commit/daf8a65bd246902bc61ca497d659a1300a13718b))
* **main:** hek3ster 1.0.1-alpha.5 ([#35](https://github.com/magenx/hek3ster/issues/35)) ([3fbb64d](https://github.com/magenx/hek3ster/commit/3fbb64d23855c5beacb4fa831d93c6ca0965421c))
* **main:** hek3ster 1.0.1-alpha.6 ([#40](https://github.com/magenx/hek3ster/issues/40)) ([33a6bc3](https://github.com/magenx/hek3ster/commit/33a6bc306b7664d7530910dc7cf69acf33dcfae8))
* **main:** hek3ster 1.0.1-alpha.7 ([#44](https://github.com/magenx/hek3ster/issues/44)) ([eeb5bbb](https://github.com/magenx/hek3ster/commit/eeb5bbb773f45e930e1a332efaa7e7e5f18a601c))
* **main:** hek3ster 1.0.1-alpha.8 ([#48](https://github.com/magenx/hek3ster/issues/48)) ([b759e20](https://github.com/magenx/hek3ster/commit/b759e2081217c92938d619d05ca61fec03a29ba5))
* **main:** hek3ster 1.0.1-alpha.9 ([#52](https://github.com/magenx/hek3ster/issues/52)) ([a12ad7c](https://github.com/magenx/hek3ster/commit/a12ad7c994f82b4fd02ab0e5416161fc1db3c604))
* Refactor binary name handling in build workflow ([#57](https://github.com/magenx/hek3ster/issues/57)) ([f6a565e](https://github.com/magenx/hek3ster/commit/f6a565e742a1948ba524f7784135e11236491555))
* Refactor build-pre-release workflow for binary naming ([#50](https://github.com/magenx/hek3ster/issues/50)) ([53140c4](https://github.com/magenx/hek3ster/commit/53140c4d996cffc277f1b44cc59819e8564c9b9a))
* Refactor build-release.yml for security and clarity ([#38](https://github.com/magenx/hek3ster/issues/38)) ([08f40b7](https://github.com/magenx/hek3ster/commit/08f40b7eeb740158e4b2c7557a8b679ec1a84d02))
* Remove commented-out code in init function ([#39](https://github.com/magenx/hek3ster/issues/39)) ([39f8cfd](https://github.com/magenx/hek3ster/commit/39f8cfdebc8b791c2d363fb1e04a4fc656e4b861))
* Remove commented-out code in init function ([#54](https://github.com/magenx/hek3ster/issues/54)) ([5a5c9c7](https://github.com/magenx/hek3ster/commit/5a5c9c76fbf45726b12910f429aca803446e750a))
* Remove macOS amd64 build target ([#21](https://github.com/magenx/hek3ster/issues/21)) ([519aa61](https://github.com/magenx/hek3ster/commit/519aa6152dbf5298068b53cc7e27d2a51abbe6d7))
* Remove old workflow logic ([a8c0aad](https://github.com/magenx/hek3ster/commit/a8c0aad69f4dbcb4fecb26b171f2c304d47ca5e1))
* Remove old workflow logic code ([#9](https://github.com/magenx/hek3ster/issues/9)) ([a8c0aad](https://github.com/magenx/hek3ster/commit/a8c0aad69f4dbcb4fecb26b171f2c304d47ca5e1))
* Remove quotes from echo statements in Makefile ([#28](https://github.com/magenx/hek3ster/issues/28)) ([972b0d2](https://github.com/magenx/hek3ster/commit/972b0d2affe39cc5d22ef280dc64a5c478dda6dd))
* Remove unnecessary comment in init function ([#47](https://github.com/magenx/hek3ster/issues/47)) ([a8dcf8b](https://github.com/magenx/hek3ster/commit/a8dcf8b3c58ac4042cfe782d739466dbedb91507))
* Remove unnecessary comment in init function ([#62](https://github.com/magenx/hek3ster/issues/62)) ([10f633e](https://github.com/magenx/hek3ster/commit/10f633eb68b004707650fb21debd43a3a6394933))
* Replace GITHUB_TOKEN with HEK3STER_TOKEN in workflow ([#46](https://github.com/magenx/hek3ster/issues/46)) ([9d10175](https://github.com/magenx/hek3ster/commit/9d10175b9f71d390ad8bbc721dad51dc0b6d7e30))
* Simplify commit message for Go formatting ([#73](https://github.com/magenx/hek3ster/issues/73)) ([404b864](https://github.com/magenx/hek3ster/commit/404b864d9df6de4d5aa3919f6eeb0759fb8226b5))
* Simplify PR title subject pattern validation ([#11](https://github.com/magenx/hek3ster/issues/11)) ([1d32ef6](https://github.com/magenx/hek3ster/commit/1d32ef6f185cd616b6246b7895d6bb785173f3b8))
* Update build-pre-release workflow for better output handling ([#61](https://github.com/magenx/hek3ster/issues/61)) ([36b1988](https://github.com/magenx/hek3ster/commit/36b1988545101eb7bca30e82008882ae4f083516))
* Update buttons to include links for actions ([#14](https://github.com/magenx/hek3ster/issues/14)) ([9c47468](https://github.com/magenx/hek3ster/commit/9c47468ce084debe0a86ad2b42a2ea25bec373d0))
* Update fetch-depth for checkout action in workflow ([#70](https://github.com/magenx/hek3ster/issues/70)) ([dfacbce](https://github.com/magenx/hek3ster/commit/dfacbce77f3fbb70dd01d59b766e1c45d84a7869))
* Update GH_TOKEN to use secrets in build workflow ([#37](https://github.com/magenx/hek3ster/issues/37)) ([faebc98](https://github.com/magenx/hek3ster/commit/faebc98d379ed542581dbd9254a2b487dca7bd6c))
* Update Makefile ([#33](https://github.com/magenx/hek3ster/issues/33)) ([44093c1](https://github.com/magenx/hek3ster/commit/44093c1e64181c67401d8da7fff5a0b9cfbfcab2))
* Update release workflow to ignore more paths ([#16](https://github.com/magenx/hek3ster/issues/16)) ([39f9a40](https://github.com/magenx/hek3ster/commit/39f9a40a34394124fad05eb78157d34d9d885c60))
* Update VERSION in Makefile to 0.0.0 ([e6aa0b9](https://github.com/magenx/hek3ster/commit/e6aa0b906cb9849d675195d16f3f058d08a2e372))
* Update VERSION in Makefile to 0000 ([#10](https://github.com/magenx/hek3ster/issues/10)) ([e6aa0b9](https://github.com/magenx/hek3ster/commit/e6aa0b906cb9849d675195d16f3f058d08a2e372))
* Update workflow for Go app release process ([#68](https://github.com/magenx/hek3ster/issues/68)) ([89bcfcf](https://github.com/magenx/hek3ster/commit/89bcfcf1fda32668c6559db4be051dd0b9a18b78))
* Use CGO conditional variable flag ([b412569](https://github.com/magenx/hek3ster/commit/b41256954ee5fa05dbb6d930abb8a072c7a730ec))
* **workflow:** Add GitHub Actions workflow for Go formatting ([#67](https://github.com/magenx/hek3ster/issues/67)) ([6c0fcd9](https://github.com/magenx/hek3ster/commit/6c0fcd9b2f0ea9a0dc114299b4074a99802226db))

## [1.0.1-alpha.13](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.12...v1.0.1-alpha.13) (2026-01-14)


### 🚦 Maintenance

* Delete .github/workflows/build-release.yml ([#66](https://github.com/magenx/hek3ster/issues/66)) ([4ef0868](https://github.com/magenx/hek3ster/commit/4ef0868b33782b443e0cc95e107281977f99c375))
* Enhance installation summary formatting in YAML ([#65](https://github.com/magenx/hek3ster/issues/65)) ([35f322c](https://github.com/magenx/hek3ster/commit/35f322c240a32c7cb0a56e3031baca4498b4b13d))
* Fix formatting issues in completion.go ([#71](https://github.com/magenx/hek3ster/issues/71)) ([4a25c47](https://github.com/magenx/hek3ster/commit/4a25c47afa43400452626de745b850ed6154378f))
* Update fetch-depth for checkout action in workflow ([#70](https://github.com/magenx/hek3ster/issues/70)) ([dfacbce](https://github.com/magenx/hek3ster/commit/dfacbce77f3fbb70dd01d59b766e1c45d84a7869))
* Update workflow for Go app release process ([#68](https://github.com/magenx/hek3ster/issues/68)) ([89bcfcf](https://github.com/magenx/hek3ster/commit/89bcfcf1fda32668c6559db4be051dd0b9a18b78))
* **workflow:** Add GitHub Actions workflow for Go formatting ([#67](https://github.com/magenx/hek3ster/issues/67)) ([6c0fcd9](https://github.com/magenx/hek3ster/commit/6c0fcd9b2f0ea9a0dc114299b4074a99802226db))

## [1.0.1-alpha.12](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.11...v1.0.1-alpha.12) (2026-01-14)


### 🚦 Maintenance

* Remove unnecessary comment in init function ([#62](https://github.com/magenx/hek3ster/issues/62)) ([10f633e](https://github.com/magenx/hek3ster/commit/10f633eb68b004707650fb21debd43a3a6394933))
* Update build-pre-release workflow for better output handling ([#61](https://github.com/magenx/hek3ster/issues/61)) ([36b1988](https://github.com/magenx/hek3ster/commit/36b1988545101eb7bca30e82008882ae4f083516))

## [1.0.1-alpha.11](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.10...v1.0.1-alpha.11) (2026-01-14)


### 🚦 Maintenance

* Add comment to init function in completion.go ([#58](https://github.com/magenx/hek3ster/issues/58)) ([742a763](https://github.com/magenx/hek3ster/commit/742a76369e0b7638ded7d35436d20c097841707b))
* Refactor binary name handling in build workflow ([#57](https://github.com/magenx/hek3ster/issues/57)) ([f6a565e](https://github.com/magenx/hek3ster/commit/f6a565e742a1948ba524f7784135e11236491555))

## [1.0.1-alpha.10](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.9...v1.0.1-alpha.10) (2026-01-13)


### 🚦 Maintenance

* Remove commented-out code in init function ([#54](https://github.com/magenx/hek3ster/issues/54)) ([5a5c9c7](https://github.com/magenx/hek3ster/commit/5a5c9c76fbf45726b12910f429aca803446e750a))

## [1.0.1-alpha.9](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.8...v1.0.1-alpha.9) (2026-01-13)


### 🚦 Maintenance

* Add comment placeholder in init function ([#51](https://github.com/magenx/hek3ster/issues/51)) ([604c36a](https://github.com/magenx/hek3ster/commit/604c36a921e549839294cd50714734f847380667))
* Refactor build-pre-release workflow for binary naming ([#50](https://github.com/magenx/hek3ster/issues/50)) ([53140c4](https://github.com/magenx/hek3ster/commit/53140c4d996cffc277f1b44cc59819e8564c9b9a))

## [1.0.1-alpha.8](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.7...v1.0.1-alpha.8) (2026-01-13)


### 🚦 Maintenance

* Remove unnecessary comment in init function ([#47](https://github.com/magenx/hek3ster/issues/47)) ([a8dcf8b](https://github.com/magenx/hek3ster/commit/a8dcf8b3c58ac4042cfe782d739466dbedb91507))
* Replace GITHUB_TOKEN with HEK3STER_TOKEN in workflow ([#46](https://github.com/magenx/hek3ster/issues/46)) ([9d10175](https://github.com/magenx/hek3ster/commit/9d10175b9f71d390ad8bbc721dad51dc0b6d7e30))

## [1.0.1-alpha.7](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.6...v1.0.1-alpha.7) (2026-01-13)


### 🚦 Maintenance

* Add comment placeholder in init function ([#43](https://github.com/magenx/hek3ster/issues/43)) ([a0e481a](https://github.com/magenx/hek3ster/commit/a0e481ac4a5febf30cb0a05ec2dcaab16edaf1ea))
* Enhance build-pre-release.yml for artifact uploads ([#42](https://github.com/magenx/hek3ster/issues/42)) ([0062aa1](https://github.com/magenx/hek3ster/commit/0062aa132ef155bad5d8238f44b833464eeb6bff))

## [1.0.1-alpha.6](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.5...v1.0.1-alpha.6) (2026-01-13)


### 🚦 Maintenance

* Refactor build-release.yml for security and clarity ([#38](https://github.com/magenx/hek3ster/issues/38)) ([08f40b7](https://github.com/magenx/hek3ster/commit/08f40b7eeb740158e4b2c7557a8b679ec1a84d02))
* Remove commented-out code in init function ([#39](https://github.com/magenx/hek3ster/issues/39)) ([39f8cfd](https://github.com/magenx/hek3ster/commit/39f8cfdebc8b791c2d363fb1e04a4fc656e4b861))
* Update GH_TOKEN to use secrets in build workflow ([#37](https://github.com/magenx/hek3ster/issues/37)) ([faebc98](https://github.com/magenx/hek3ster/commit/faebc98d379ed542581dbd9254a2b487dca7bd6c))

## [1.0.1-alpha.5](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.4...v1.0.1-alpha.5) (2026-01-13)


### 🚦 Maintenance

* Add comment placeholder in init function ([#34](https://github.com/magenx/hek3ster/issues/34)) ([98803ca](https://github.com/magenx/hek3ster/commit/98803ca0113650164441b39c31308834e9a4ac10))
* Add permissions for contents in build workflow ([#32](https://github.com/magenx/hek3ster/issues/32)) ([c25f1c2](https://github.com/magenx/hek3ster/commit/c25f1c273b635d99e496c2d5f07f45abfd1148a7))
* Enhance build script with binary checks and uploads ([#31](https://github.com/magenx/hek3ster/issues/31)) ([c014b2b](https://github.com/magenx/hek3ster/commit/c014b2ba84ed3971abf55299e18144a37d816249))
* Update Makefile ([#33](https://github.com/magenx/hek3ster/issues/33)) ([44093c1](https://github.com/magenx/hek3ster/commit/44093c1e64181c67401d8da7fff5a0b9cfbfcab2))

## [1.0.1-alpha.4](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.3...v1.0.1-alpha.4) (2026-01-13)


### 🚦 Maintenance

* Remove quotes from echo statements in Makefile ([#28](https://github.com/magenx/hek3ster/issues/28)) ([972b0d2](https://github.com/magenx/hek3ster/commit/972b0d2affe39cc5d22ef280dc64a5c478dda6dd))

## [1.0.1-alpha.3](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.2...v1.0.1-alpha.3) (2026-01-13)


### 🚦 Maintenance

* Change cancel-in-progress to false in workflow ([#24](https://github.com/magenx/hek3ster/issues/24)) ([fb7dc90](https://github.com/magenx/hek3ster/commit/fb7dc90fd762403aa346af410c49e79161e42ae1))
* Fix build-all target in Makefile ([#25](https://github.com/magenx/hek3ster/issues/25)) ([d6a7ce4](https://github.com/magenx/hek3ster/commit/d6a7ce4802cc4fab5132499b179d503b452a0592))

## [1.0.1-alpha.2](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha.1...v1.0.1-alpha.2) (2026-01-13)


### 🚦 Maintenance

* Add GH_TOKEN to build environment ([#18](https://github.com/magenx/hek3ster/issues/18)) ([d4a31d9](https://github.com/magenx/hek3ster/commit/d4a31d9c09acd5fd63227a945697b56e1ffd9737))
* Add GH_TOKEN to build-release workflow ([#19](https://github.com/magenx/hek3ster/issues/19)) ([a6547c9](https://github.com/magenx/hek3ster/commit/a6547c9a2522b65fd083c3c40605ecf5982dd442))
* Delete CODE_OF_CONDUCT.md ([#20](https://github.com/magenx/hek3ster/issues/20)) ([36d35c5](https://github.com/magenx/hek3ster/commit/36d35c5e4a8e9d087f743a2ef19ad7f4ef503e9d))
* Exclude some files from workflow  trigger ([2bd0504](https://github.com/magenx/hek3ster/commit/2bd050422a75b4f6af32e8ed910358e7fa16e15c))
* Remove macOS amd64 build target ([#21](https://github.com/magenx/hek3ster/issues/21)) ([519aa61](https://github.com/magenx/hek3ster/commit/519aa6152dbf5298068b53cc7e27d2a51abbe6d7))

## [1.0.1-alpha.1](https://github.com/magenx/hek3ster/compare/v1.0.1-alpha...v1.0.1-alpha.1) (2026-01-13)


### 🚦 Maintenance

* Update buttons to include links for actions ([#14](https://github.com/magenx/hek3ster/issues/14)) ([9c47468](https://github.com/magenx/hek3ster/commit/9c47468ce084debe0a86ad2b42a2ea25bec373d0))
* Use CGO conditional variable flag ([b412569](https://github.com/magenx/hek3ster/commit/b41256954ee5fa05dbb6d930abb8a072c7a730ec))

## [1.0.1-alpha](https://github.com/magenx/hek3ster/compare/v1.0.0...v1.0.1-alpha) (2026-01-13)


### 🚦 Maintenance

* Build workflow update ([da43e6c](https://github.com/magenx/hek3ster/commit/da43e6c6de9d7cdd0caff0e4910fbaa89097918c))
* Build workflow update release pre release ([2c1cd97](https://github.com/magenx/hek3ster/commit/2c1cd975d386fef52ceaa0b42fa85ad95766b412))
* Remove old workflow logic ([a8c0aad](https://github.com/magenx/hek3ster/commit/a8c0aad69f4dbcb4fecb26b171f2c304d47ca5e1))
* Remove old workflow logic code ([#9](https://github.com/magenx/hek3ster/issues/9)) ([a8c0aad](https://github.com/magenx/hek3ster/commit/a8c0aad69f4dbcb4fecb26b171f2c304d47ca5e1))
* Simplify PR title subject pattern validation ([#11](https://github.com/magenx/hek3ster/issues/11)) ([1d32ef6](https://github.com/magenx/hek3ster/commit/1d32ef6f185cd616b6246b7895d6bb785173f3b8))
* Update VERSION in Makefile to 0.0.0 ([e6aa0b9](https://github.com/magenx/hek3ster/commit/e6aa0b906cb9849d675195d16f3f058d08a2e372))
* Update VERSION in Makefile to 0000 ([#10](https://github.com/magenx/hek3ster/issues/10)) ([e6aa0b9](https://github.com/magenx/hek3ster/commit/e6aa0b906cb9849d675195d16f3f058d08a2e372))
