# Changelog

## [4.2.0](https://github.com/andrii-yeremenko/news-aggregator/compare/news-aggregator@v4.1.0...news-aggregator@v4.2.0) (2024-10-04)


### Features

* add GitOps `news-aggregator` app. ([3d62eb6](https://github.com/andrii-yeremenko/news-aggregator/commit/3d62eb6bff1bbf400b18cc2739931787964ade08))


### Bug Fixes

* add to project destinations `kube-system` namespace. ([2c9198a](https://github.com/andrii-yeremenko/news-aggregator/commit/2c9198ac8b9849cf07b47ea993af9cb3a6c0e787))
* add to project destinations `kube-system` namespace. ([3289e37](https://github.com/andrii-yeremenko/news-aggregator/commit/3289e375953a901fb332ced5526f3342ee158230))


### Reverts

* unnecessary changes. ([84f2092](https://github.com/andrii-yeremenko/news-aggregator/commit/84f209251c55cc8ba64530167d44a41a8c9a3f33))
* unnecessary changes. ([3c0eb88](https://github.com/andrii-yeremenko/news-aggregator/commit/3c0eb88753072d2c54811bcfd636dcbaa8eb18e0))

## [4.1.0](https://github.com/andrii-yeremenko/news-aggregator/compare/news-aggregator@v4.0.0...news-aggregator@v4.1.0) (2024-09-29)


### Features

* add helm check if AWS `secretKey` and `accessKey` are provided. ([3dc8ac0](https://github.com/andrii-yeremenko/news-aggregator/commit/3dc8ac0b95cf63bfb307e69853d5e52df872de14))
* inject all secret keys on start of operator and chart. ([e2d31fc](https://github.com/andrii-yeremenko/news-aggregator/commit/e2d31fcadd893db686eff4b56e824e9c5a60b72b))


### Bug Fixes

* nested Taskfile variables conflict. ([348f7e0](https://github.com/andrii-yeremenko/news-aggregator/commit/348f7e0cd9c1fd06e3e1e52b13815614241362e5))
* operator and updater Taskfiles path problem. ([2498699](https://github.com/andrii-yeremenko/news-aggregator/commit/2498699403ba27126080c918e6c190aa9d7ad7ad))
* operator and updater Taskfiles path problem. ([4398753](https://github.com/andrii-yeremenko/news-aggregator/commit/43987533d06a170ca78959d48b6b4e889a1a62c0))

## [4.0.0](https://github.com/andrii-yeremenko/news-aggregator/compare/news-aggregator@v3.0.0...news-aggregator@v4.0.0) (2024-09-26)


### ⚠ BREAKING CHANGES

* update news-aggregator server's helm chart version from `0.1.0` to `1.0.0`.
* update https news-aggregator server's image version from `1.0.0` to `2.0.0`.

### Miscellaneous Chores

* update https news-aggregator server's image version from `1.0.0` to `2.0.0`. ([73002bf](https://github.com/andrii-yeremenko/news-aggregator/commit/73002bf221db5cd6a42f3b3e88033067e99e99cb))
* update news-aggregator server's helm chart version from `0.1.0` to `1.0.0`. ([61f8b93](https://github.com/andrii-yeremenko/news-aggregator/commit/61f8b93ea04fd1b64c1fd3127db4cd6c1ae358ec))

## [3.0.0](https://github.com/andrii-yeremenko/news-aggregator/compare/news-aggregator@v2.0.0...news-aggregator@v3.0.0) (2024-09-26)


### ⚠ BREAKING CHANGES

* improve error description in `HotNews` webhook.
* add more logging in `Reconcile()` function.

### Features

* add `cert-manager`. ([f329d90](https://github.com/andrii-yeremenko/news-aggregator/commit/f329d9045dc515cbef6aefe0d60c9fb5456b979c))
* add `HotNews` CRD and controller. ([2769dfc](https://github.com/andrii-yeremenko/news-aggregator/commit/2769dfc7663b370f3fd21718933391dac456e9cd))
* add additional checks in `setup-feeds` task. ([6cf8e26](https://github.com/andrii-yeremenko/news-aggregator/commit/6cf8e26e1716583ae241aa75dddd74c684656b0a))
* add Conditions to HotNews. ([f8e987c](https://github.com/andrii-yeremenko/news-aggregator/commit/f8e987cc0469e85750da271c76e8265a0e2693b0))
* add delete event handling. ([f59d5be](https://github.com/andrii-yeremenko/news-aggregator/commit/f59d5bebdd7e0a4947461cf997bb0067e3a5efc1))
* add Feed `Failed` condition type. ([128ff98](https://github.com/andrii-yeremenko/news-aggregator/commit/128ff98d03d9744e1489b0a4149c440bed1e8a22))
* add Feed fields validation. ([35b24e3](https://github.com/andrii-yeremenko/news-aggregator/commit/35b24e3c6914a1fef9bd734ea4cbef41a1235204))
* add feeds validation, remove hardcoded variables `configMapName` and `configMapNamespace`. ([193dc0e](https://github.com/andrii-yeremenko/news-aggregator/commit/193dc0ed7a9c2e73b6d20f842a6b6061d89664cb))
* add finalizer and manage OwnerReferences in HotNewsReconciler ([99f5594](https://github.com/andrii-yeremenko/news-aggregator/commit/99f55948985aa43a04f52059b9274e04cbd6adc2))
* add more logging in `Reconcile()` function. ([1462b60](https://github.com/andrii-yeremenko/news-aggregator/commit/1462b60740333eddd0e4ad5a9241ed1a84b44377))
* add to feed's webhook `ValidateDelete` ([cd7b8f4](https://github.com/andrii-yeremenko/news-aggregator/commit/cd7b8f429bf52bdaabed2193a9c9dd4a495ba3d1))
* add validation webhook for `feed-groups` ConfigMap. ([b31ab33](https://github.com/andrii-yeremenko/news-aggregator/commit/b31ab33ac86f846e18c1b197eecdbb2727b2dd88))
* implement `Feed` controller. ([178ef24](https://github.com/andrii-yeremenko/news-aggregator/commit/178ef244f74a3d2b86d964c5e93e992d1fbef0f8))
* implement validation webhook for `HotNews` CRD. ([a7542bc](https://github.com/andrii-yeremenko/news-aggregator/commit/a7542bcf41f008e9d0546369b167b5f1dff45c82))
* improve error description in `HotNews` webhook. ([e70ee71](https://github.com/andrii-yeremenko/news-aggregator/commit/e70ee71a39057427df90c46a24c595566e686c59))
* introduce custom predicates to watches. ([86404be](https://github.com/andrii-yeremenko/news-aggregator/commit/86404be6f37e027b80bc7e68025765d3e0d530a0))


### Bug Fixes

* `HotNews` reconciliation works even if no updated `ConfigMap` is used. ([c045136](https://github.com/andrii-yeremenko/news-aggregator/commit/c0451369da661c15e1c5f9be774bfd65075de6b1))
* bug with `OwnerReference`'s. ([2f6bfe4](https://github.com/andrii-yeremenko/news-aggregator/commit/2f6bfe4f34b4d50c1f08f4061c73a1bc0267b486))
* bug with HotNews reconciliation. ([3925bd6](https://github.com/andrii-yeremenko/news-aggregator/commit/3925bd6f3b9c38649f06a501ccf340cd50d86c7f))
* call reconcile only if `HotNews` contains target `Feed`. ([65f0e86](https://github.com/andrii-yeremenko/news-aggregator/commit/65f0e867200d9768ec39553d9928abcc2e14aee0))
* correct context.WithTimeout() usage. ([7e45b30](https://github.com/andrii-yeremenko/news-aggregator/commit/7e45b3081440850aaf15781b7a0da1aa1e1b711f))
* correct context.WithTimeout() usage. ([146ef14](https://github.com/andrii-yeremenko/news-aggregator/commit/146ef149e575647077e92b90be703f442b0bf9d9))
* fix `HotNews` controller's watch events to handle feed and configmap updation. ([c1775b5](https://github.com/andrii-yeremenko/news-aggregator/commit/c1775b56dab4f7217faa6dcd1769b10c11cc6c8c))
* make webhook check Feed `Spec.Name` instead `Name`. ([f590437](https://github.com/andrii-yeremenko/news-aggregator/commit/f590437fe6e192edd3d268c4bab53269d16a9286))
* remove creation of new `context`. ([c7e1243](https://github.com/andrii-yeremenko/news-aggregator/commit/c7e1243e638c212ebd5b5930781702f8551272d3))
* remove redundant code. ([21fcbdb](https://github.com/andrii-yeremenko/news-aggregator/commit/21fcbdb693911ff0d56f65d962085d8a8fad7ae5))
* remove redundant validations that already cover webhook. ([e5b3254](https://github.com/andrii-yeremenko/news-aggregator/commit/e5b32544cd6356b0a5d8183733be309b061e9451))
* remove suite test for webhooks. ([ea51232](https://github.com/andrii-yeremenko/news-aggregator/commit/ea512321fdaa6c66c91e47d2b897dac926e66d83))
* update status bug on deletion. ([f8ef4ee](https://github.com/andrii-yeremenko/news-aggregator/commit/f8ef4eebd19f45afacb470a7ce9a4553c74b341b))
* webhook annotation ([2f689dc](https://github.com/andrii-yeremenko/news-aggregator/commit/2f689dcbaa2e03e7eabad9d3d873d76d6de84d4e))
