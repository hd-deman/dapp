---
title: Using artifacts
sidebar: reference
permalink: reference/build/artifact.html
author: Alexey Igrychev <alexey.igrychev@flant.com>
summary: |
  <a href="https://docs.google.com/drawings/d/e/2PACX-1vRPnqkxbv8wSziAE7QVhcP4rsb58AfIGOmOvVUbWKtZdvNhGItnL0RX8ZFZgCxxNZTtYdZ6YbVuItix/pub?w=1914&amp;h=721" data-featherlight="image">
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vRPnqkxbv8wSziAE7QVhcP4rsb58AfIGOmOvVUbWKtZdvNhGItnL0RX8ZFZgCxxNZTtYdZ6YbVuItix/pub?w=957&amp;h=360">
  </a>
---

The size of the image can be increased several times due to the assembly tools and source files, while the user does not need them. 

To solve such problems, the Docker community suggests doing the installation of tools, the assembly, and removal of tools in one step.

```
RUN “download-source && cmd && cmd2 && remove-source”
```

> To obtain a similar effect when using dapp, it is sufficient to describe the instructions in one _user stage_. For example, the _shell assembly instructions_ for the _install stage_ (similarly for _ansible_):
```yaml
shell:
  install:
  - "download-source"
  - "cmd"
  - "cmd2"
  - "remove-source"
```

However, with this method, it isn't possible to use caching and need to install toolkit regularly.

Another solution is to use multi-stage builds, which are supported starting with Docker 17.05.

```
FROM node:latest AS storefront
WORKDIR /usr/src/atsea/app/react-app
COPY react-app .
RUN npm install
RUN npm run build

FROM maven:latest AS appserver
WORKDIR /usr/src/atsea
COPY pom.xml .
RUN mvn -B -f pom.xml -s /usr/share/maven/ref/settings-docker.xml dependency:resolve
COPY . .
RUN mvn -B -s /usr/share/maven/ref/settings-docker.xml package -DskipTests

FROM java:8-jdk-alpine
RUN adduser -Dh /home/gordon gordon
WORKDIR /static
COPY --from=storefront /usr/src/atsea/app/react-app/build/ .
WORKDIR /app
COPY --from=appserver /usr/src/atsea/target/AtSea-0.0.1-SNAPSHOT.jar .
ENTRYPOINT ["java", "-jar", "/app/AtSea-0.0.1-SNAPSHOT.jar"]
CMD ["--spring.profiles.active=postgres"]
```

The meaning of such an approach is as follows, describe several auxiliary images, and selectively copy artifacts from one image to another, leaving behind everything you don’t want in the final image.

Dapp suggests an alternative in the form of _artifact dimg_, which are built according to the same rules as _dimg_, but with slight changes in the _stage conveyor_.

> Why doesn't dapp use multi-stage? Historically, _artifacts_ appeared much earlier than Docker multi-stage, and dapp approach gives more flexibility when working with auxiliary images.

## What is an artifact?

***Artifact*** is special dimg that is used by _dimgs_ and _artifacts_ to isolate the build process and build tools resources (environments, software, data).
          
_Artifact_ cannot be [tagged like _dimg_]({{ site.baseurl }}/reference/registry/image_naming.html#dapp-tag-procedure) and used as standalone application.

Using artifacts, you can independently assemble an unlimited number of components, and also solving the following problems:

- The application can consist of a set of components, and each has its dependencies. With a standard assembly, you should rebuild all every time, but you want to assemble each one on-demand.
- Components need to be assembled in other environments.

Importing _artifacts resources_ are described at destination _dimg_ or _artifact_ by `import` directive records.

## Configuration

The _dappfile configuration_ of the _artifact_ is not much different from the configuration of _dimg_. Each _artifact_ should be described in a separate YAML document.

The instructions associated with the _from stage_, namely the [_base image_]({{ site.baseurl }}/reference/build/base_image.html) and [mounts]({{ site.baseurl }}/reference/build/mount_directive.html), remain unchanged.

The _docker_instructions stage_ and the [corresponding instructions]({{ site.baseurl }}/reference/build/docker_directive.html) are not supported for the _artifact_. A _artifact_ is an assembly tool and only the data stored in it is required.

The remaining _stages_ and instructions are considered further separately.

### Naming

<div class="summary" markdown="1">
```yaml
artifact: <artifact name>
```
</div>

_Artifact images_ are declared with `artifact` directive: `artifact: <artifact name>`. Unlike the [naming of the _dimg_]({{ site.baseurl }}/reference/build/naming.html), the artifact has no limitations associated with docker naming convention, as used only internal.

```yaml
artifact: "application assets"
```  

The _artifact name_ is used to specify the artifact in the [_artifact resources import_ description](#importing-artifacts) of the _dimg_ or _artifact_.

### Adding source code from git repositories

<div class="summary" markdown="1">

<a href="https://docs.google.com/drawings/d/e/2PACX-1vRYyGoELol94us3ahKhn_I8efTOajXntgf-WhB6QPykKZrdm296B6TJm3YAE-DTmkBpKTm9AvZQbZC5/pub?w=2031&amp;h=144" data-featherlight="image">
<img src="https://docs.google.com/drawings/d/e/2PACX-1vRYyGoELol94us3ahKhn_I8efTOajXntgf-WhB6QPykKZrdm296B6TJm3YAE-DTmkBpKTm9AvZQbZC5/pub?w=1016&amp;h=72">
</a>

```yaml
git:
- ...
  stageDependencies:
    ...
    buildArtifact: 
    - <relative_path or glob>
```

</div>

Unlike with _dimg_, _artifact stage conveyor_ has no _g_a_cache_ and _g_a_latest_patch_ stages.

> Dapp implements optional dependence on changes in git repositories for _artifacts_. Thus, by default dapp ignores them and _artifact image_ is cached after the first assembly, but you can specify any dependencies for assembly instructions.

Read about working with _git repositories_ in the corresponding [article]({{ site.baseurl }}/reference/build/git_directive.html).

### Running assembly instructions

<div class="summary" markdown="1">

<a href="https://docs.google.com/drawings/d/e/2PACX-1vSEje1gsyjI89m4lh6PqDEFcwa7NsLeTnbju1hZ7G4AJ2S4f_nJlczEne6rbpuvtoDkbBCqhu-i5dnT/pub?w=2031&amp;h=144" data-featherlight="image">
<img src="https://docs.google.com/drawings/d/e/2PACX-1vSEje1gsyjI89m4lh6PqDEFcwa7NsLeTnbju1hZ7G4AJ2S4f_nJlczEne6rbpuvtoDkbBCqhu-i5dnT/pub?w=1016&amp;h=72">
</a>

<div class="tab">
    <button class="tablinks active" onclick="openTab(event, 'shell')">Shell</button>
    <button class="tablinks" onclick="openTab(event, 'ansible')">Ansible</button>
</div>

<div id="shell" class="tabcontent active" markdown="1">
```yaml
shell:
  ...
  buildArtifact:
  - <bash command>
  buildArtifactCacheVersion: <arbitrary string>
```
</div>

<div id="ansible" class="tabcontent" markdown="1">
```yaml
ansible:
  ...
  buildArtifact:
  - <task>
  buildArtifactCacheVersion: <arbitrary string>
```
</div>

</div>

When describing an _artifact_, an additional _build_artifact user stage_ is available, which is no different from the _install_, _before_setup_, and _setup_ stages. As well as the listed _stages_, _build_artifact_ also has a dependent git stage.

<a href="https://docs.google.com/drawings/d/e/2PACX-1vTd0XO1HQdwWQRB-QCNmxhJcaBdZG5m4YktzhMXLB4hdu8NxEEnWZvbaivwK13pEfddxtiHNXzgjhal/pub?w=1917&amp;h=432" data-featherlight="image">
<img src="https://docs.google.com/drawings/d/e/2PACX-1vTd0XO1HQdwWQRB-QCNmxhJcaBdZG5m4YktzhMXLB4hdu8NxEEnWZvbaivwK13pEfddxtiHNXzgjhal/pub?w=959&amp;h=216">
</a>

If there are no dependencies on files specified in git `stageDependencies` directive for _artifact user stages_, the image is cached after the first build and will no longer be reassembled while the _stages cache_ exists.

> If the artifact should be rebuilt on any change in the related git repository, you should specify the _stageDependency_ `**/*` for any _user stage_, e.g., for _build_artifact stage_:
```yaml
git:
- to: /
  stageDependencies:
    buildArtifact: "**/*"
```

Read about working with _assembly instructions_ in the corresponding [article]({{ site.baseurl }}/reference/build/assembly_instructions.html).

### Importing artifacts

{% include_relative import_artifacts_partial.md %}

## All directives
```yaml
artifact: <artifact_name>
from: <image>
fromCacheVersion: <version>
fromDimg: <dimg_name>
fromDimgArtifact: <artifact_name>
git:
# local git
- as: <custom_name>
  add: <absolute_path>
  to: <absolute_path>
  owner: <owner>
  group: <group>
  includePaths:
  - <relative_path or glob>
  excludePaths:
  - <relative_path or glob>
  stageDependencies:
    install:
    - <relative_path or glob>
    beforeSetup:
    - <relative_path or glob>
    setup:
    - <relative_path or glob>
    buildArtifact:
    - <relative_path or glob>
# remote git
- url: <git_repo_url>
  branch: <branch_name>
  commit: <commit>
  tag: <tag>
  as: <custom_name>
  add: <absolute_path>
  to: <absolute_path>
  owner: <owner>
  group: <group>
  includePaths:
  - <relative_path or glob>
  excludePaths:
  - <relative_path or glob>
  stageDependencies:
    install:
    - <relative_path or glob>
    beforeSetup:
    - <relative_path or glob>
    setup:
    - <relative_path or glob>
    buildArtifact:
    - <relative_path or glob>
shell:
  beforeInstall:
  - <cmd>
  install:
  - <cmd>
  beforeSetup:
  - <cmd>
  setup:
  - <cmd>
  buildArtifact:
  - <cmd>
  cacheVersion: <version>
  beforeInstallCacheVersion: <version>
  installCacheVersion: <version>
  beforeSetupCacheVersion: <version>
  setupCacheVersion: <version>
  buildArtifactCacheVersion: <version>
ansible:
  beforeInstall:
  - <task>
  install:
  - <task>
  beforeSetup:
  - <task>
  setup:
  - <task>
  buildArtifact:
  - <task>
  cacheVersion: <version>
  beforeInstallCacheVersion: <version>
  installCacheVersion: <version>
  beforeSetupCacheVersion: <version>
  setupCacheVersion: <version>
  buildArtifactCacheVersion: <version>
mount:
- from: build_dir
  to: <absolute_path>
- from: tmp_dir
  to: <absolute_path>
- fromPath: <absolute_path>
  to: <absolute_path>
import:
- artifact: <artifact name>
  before: <install || setup>
  after: <install || setup>
  add: <absolute path>
  to: <absolute path>
  owner: <owner>
  group: <group>
  includePaths:
  - <relative path or glob>
  excludePaths:
  - <relative path or glob>
asLayers: <false || true>
```
