#!/usr/bin/env python3
"""Check relative links and Sphinx references in the documentation tree.

No third-party dependencies. By default it checks the Chinese Sphinx
sources (docs/cn/source). Use --scope all to also check AGENTS.md,
README.md, features/, docs/agents/ and docs/reference/.

Exit code: 0 when every target exists, 1 when at least one is missing.
"""
import argparse
import re
import sys
import urllib.parse
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
LANG_SOURCE = {
    "cn": REPO_ROOT / "docs" / "cn" / "source",
    "en": REPO_ROOT / "docs" / "en" / "source",
}
AGENT_PATHS = [
    REPO_ROOT / "AGENTS.md",
    REPO_ROOT / "README.md",
    REPO_ROOT / "features",
    REPO_ROOT / "docs" / "agents",
    REPO_ROOT / "docs" / "reference",
]

MD_LINK_RE = re.compile(r"!?\[[^\]]*\]\(([^)\s]+)(?:\s+['\"][^'\"]*['\"])?\)")
RST_DOC_RE = re.compile(r":doc:`([^`]+)`")
RST_IMAGE_RE = re.compile(r"^\s*\.\. (?:image|figure)::\s*(\S+)")
HTML_IMG_RE = re.compile(r"<img[^>]+src=['\"]([^'\"]+)['\"]", re.IGNORECASE)
TOCTREE_RE = re.compile(r"^\s*\.\. toctree::")
EXTERNAL_RE = re.compile(r"^(?:[a-zA-Z][a-zA-Z0-9+.-]*:|//|#|$)")


def iter_visible_lines(path):
    """Yield (lineno, line) outside fenced code blocks."""
    in_fence = False
    fence = ""
    text = path.read_text(encoding="utf-8", errors="replace")
    for lineno, line in enumerate(text.splitlines(), 1):
        stripped = line.lstrip()
        if stripped.startswith("```") or stripped.startswith("~~~"):
            marker = stripped[:3]
            # A line with two markers is an inline code span, not a fence.
            if stripped.count(marker) >= 2:
                continue
            if not in_fence:
                in_fence = True
                fence = marker
            elif marker == fence:
                in_fence = False
                fence = ""
            continue
        if not in_fence:
            yield lineno, line


def iter_toctree_lines(path):
    """Yield (lineno, target) for entries under .. toctree::."""
    in_tree = False
    for lineno, line in iter_visible_lines(path):
        if TOCTREE_RE.match(line):
            in_tree = True
            continue
        if not in_tree:
            continue
        if not line.strip():
            continue
        if line.lstrip().startswith(":"):
            continue
        if line.startswith((" ", "\t")):
            target = line.strip()
            if target and not target.startswith(".."):
                yield lineno, target
        else:
            in_tree = False


def normalize_doc_target(raw):
    """Extract the path from :doc:`Title <path>` or :doc:`path`."""
    raw = raw.strip()
    match = re.match(r"^(.*?)\s*<([^>]+)>\s*$", raw)
    if match:
        return match.group(2).strip()
    return raw


def iter_targets(path):
    """Yield (kind, target, lineno) for links and directives in a file."""
    for lineno, line in iter_visible_lines(path):
        for match in MD_LINK_RE.finditer(line):
            yield "md", match.group(1), lineno
        for match in RST_DOC_RE.finditer(line):
            yield "doc", normalize_doc_target(match.group(1)), lineno
        for match in RST_IMAGE_RE.finditer(line):
            yield "md", match.group(1), lineno
        for match in HTML_IMG_RE.finditer(line):
            yield "md", match.group(1), lineno


def should_skip(target):
    if not target:
        return True
    if EXTERNAL_RE.match(target):
        return True
    if "..." in target or "${" in target or "{{" in target or "<" in target:
        return True
    return False


def candidates_for(path, kind, target, source_root):
    target = urllib.parse.unquote(target)
    target = target.split("#", 1)[0].split("?", 1)[0]
    if not target:
        return []
    if target.startswith("/"):
        base = source_root
        target = target.lstrip("/")
    else:
        base = path.parent
    candidates = [base / target]
    if kind == "doc" or not Path(target).suffix:
        candidates.append(base / (target + ".md"))
        candidates.append(base / (target + ".rst"))
    return candidates


def source_root_for(path):
    for root in LANG_SOURCE.values():
        try:
            path.relative_to(root)
            return root
        except ValueError:
            pass
    return REPO_ROOT


def check_file(path, missing, syntax):
    source_root = source_root_for(path)
    for kind, target, lineno in iter_targets(path):
        if path.suffix == ".md" and kind == "doc":
            syntax.append((path, lineno, target))
        if should_skip(target):
            continue
        candidates = candidates_for(path, kind, target, source_root)
        if any(candidate.exists() for candidate in candidates):
            continue
        missing.append((path, lineno, target, candidates[0] if candidates else path))
    for lineno, target in iter_toctree_lines(path):
        if should_skip(target):
            continue
        candidates = candidates_for(path, "doc", target, source_root)
        if any(candidate.exists() for candidate in candidates):
            continue
        missing.append((path, lineno, target, candidates[0]))


def iter_files(lang, scope):
    files = []
    if scope in ("sphinx", "all"):
        root = LANG_SOURCE[lang]
        if root.is_dir():
            files.extend(sorted(root.rglob("*.md")))
            files.extend(sorted(root.rglob("*.rst")))
    if scope in ("agent", "all"):
        for path in AGENT_PATHS:
            if path.is_file() and path.suffix in (".md", ".rst"):
                files.append(path)
            elif path.is_dir():
                files.extend(sorted(path.rglob("*.md")))
    return list(dict.fromkeys(files))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--lang", choices=["cn", "en", "all"], default="cn")
    parser.add_argument("--scope", choices=["sphinx", "agent", "all"], default="sphinx")
    args = parser.parse_args()

    langs = ["cn", "en"] if args.lang == "all" else [args.lang]
    files = []
    for lang in langs:
        files.extend(iter_files(lang, args.scope))
    files = list(dict.fromkeys(files))

    missing = []
    syntax = []
    for path in files:
        check_file(path, missing, syntax)

    for path, lineno, target, resolved in missing:
        try:
            rel = path.relative_to(REPO_ROOT)
        except ValueError:
            rel = path
        print("%s:%d: %s  ->  %s" % (rel, lineno, target, resolved))

    for path, lineno, target in syntax:
        try:
            rel = path.relative_to(REPO_ROOT)
        except ValueError:
            rel = path
        print("SYNTAX %s:%d: :doc: in Markdown, use {doc}  ->  %s" % (rel, lineno, target))

    print("checked %d file(s); missing %d link target(s); markdown :doc: %d"
          % (len(files), len(missing), len(syntax)))
    return 1 if missing or syntax else 0


if __name__ == "__main__":
    sys.exit(main())
