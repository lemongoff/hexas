# Third-Party Notices

Hexas contains or adapts material from the projects listed below. This file records provenance and license boundaries; it does not replace the applicable license texts.

## go-zero

- Project: `github.com/zeromicro/go-zero`
- Version: `v1.10.3`
- Commit: `925f8a2bcc159eaf3b1da0f5fc695beac26e15ff`
- License: MIT
- Copyright notice: `Copyright (c) 2022 zeromicro`
- Included material: the framework and `tools/goctl` source baseline, subsequently modified by Hexas.
- License text: [LICENSE](LICENSE) and [LICENSES/MIT-zeromicro.txt](LICENSES/MIT-zeromicro.txt)

Hexas uses its own module/import path, `github.com/lemongoff/hexas`. The upstream path is retained only in provenance records and fixed upstream snapshots. This does not imply that Hexas is published, maintained, sponsored or endorsed by zeromicro.

## zero-skills

- Project: `github.com/zeromicro/zero-skills`
- Commit: `943a13c5d82cbd3d8f896134ea6b34182bb4c32f`
- Upstream license declaration: MIT
- Included material: selected original reference material under `ai/skills/zero-skills/upstream/`, plus a Hexas-specific skill entry.
- Adaptation record: [ai/skills/zero-skills/UPSTREAM.md](ai/skills/zero-skills/UPSTREAM.md)
- License text used for the declared MIT terms: [LICENSES/MIT-zeromicro.txt](LICENSES/MIT-zeromicro.txt)

The fixed upstream commit does not contain a standalone `LICENSE` file. Its README and skill metadata declare the project to be MIT licensed; this limitation is recorded rather than inventing a separate upstream copyright notice.

## ai-context

- Project: `github.com/zeromicro/ai-context`
- Commit: `bc525eedc924fe53b5d26f27e41595cfdb347477`
- Upstream license declaration: MIT
- Included material: workflow structure and topics rewritten for Hexas under `ai/context/`.
- Adaptation record: [ai/context/UPSTREAM.md](ai/context/UPSTREAM.md)
- License text used for the declared MIT terms: [LICENSES/MIT-zeromicro.txt](LICENSES/MIT-zeromicro.txt)

The fixed upstream commit does not contain a standalone `LICENSE` file. Its README declares the project to be MIT licensed; this limitation is recorded rather than inventing a separate upstream copyright notice.

## Contributor Covenant

- Material: [code-of-conduct.md](code-of-conduct.md)
- Source: Contributor Covenant, version 2.1
- License: Creative Commons Attribution 4.0 International (`CC-BY-4.0`)
- Modifications: translated, reorganized and adapted for Hexas, including project-specific scope, reporting, enforcement and review procedures.
- License text: [LICENSES/CC-BY-4.0.txt](LICENSES/CC-BY-4.0.txt)

## Other dependencies

Go module dependencies retain their own licenses and copyright notices. This notice is not a complete software bill of materials and does not override dependency source distributions or module metadata.
