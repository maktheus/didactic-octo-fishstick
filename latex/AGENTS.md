# Repository Guidelines

## Project Structure & Module Organization
`icomp-modelomonografia-pf.tex` is the single entrypoint and sequentially `\include`s every chapter directory, so keep the numbering intact when adding material. `meta.tex` centralizes packages, typography, and front-matter metadata, while numbered folders (`0_pre/`, `1_introducao/`, `2_fundamentos/` … `6_consideracoesFinais/`) hold the textual chapters that map to the monograph. Post-textual content lives under `8_apendice/` and `9_anexos/`, glossary acronyms stay in `0_pre/glossario.tex`, and all references belong in `referencias.bib`.

## Build, Test, and Development Commands
- `latexmk -pdf -interaction=nonstopmode icomp-modelomonografia-pf.tex` — canonical build that resolves cross-references, glossary, and bibliography.
- `latexmk -pvc -pdf icomp-modelomonografia-pf.tex` — watch mode for live previews; stop with `Ctrl+C`.
- `latexmk -c` — drop auxiliary files when switching branches or packaging artifacts.
Without `latexmk`, run `pdflatex icomp-modelomonografia-pf.tex` → `bibtex icomp-modelomonografia-pf` → `pdflatex` twice.

## Coding Style & Naming Conventions
Mirror the existing memoir/abnTeX2 style: tabs (or aligned four spaces) inside environments, blank lines between logical blocks, and `%` comments documenting template placeholders. Keep labels prefixed by type (`\label{sec:introducao}`, `\label{fig:arquitetura}`) and reuse the chapter folders’ numeric prefixes when creating new modules (`7_metodologia/cap7.tex`, for example). Place figures near their references and include them with relative paths so `abntex2` resolves them without extra TEXINPUTS tweaks.

## Testing Guidelines
Treat `latexmk -pdf -halt-on-error icomp-modelomonografia-pf.tex` as the regression suite; the build must end clean, without “Label(s) may have changed” or citation warnings. When edits touch `referencias.bib`, run one extra pass (or `bibtex`) before committing. Confirm that the PDF reflects glossary updates from `0_pre/glossario`, and that chapter numbering stays sequential.

## Commit & Pull Request Guidelines
Follow the repository-wide conventional commit pattern (`feat:`, `fix:`, `docs:`) and keep each change scoped to one chapter or asset folder so diffs stay reviewable. Note the build command you ran, attach the resulting PDF when layout changes are visible, call out any packages added to `meta.tex`, and request a reviewer tied to the affected section.

## Metadata & Configuration Tips
Update `meta.tex` whenever authorship, advisors, or PDF metadata change; it is the single source of truth for titles, course names, and chapter styling. Align glossary keys in `0_pre/glossario.tex` with their first use, keep `referencias.bib` alphabetical, and avoid touching `abntex2/*.cls` unless a documented layout change is unavoidable.
