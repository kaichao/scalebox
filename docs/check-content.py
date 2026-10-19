#!/usr/bin/env python3
"""Report coverage between Sphinx reference chapters and generated reference docs.

This is an informational baseline: it lists facts that exist in the Sphinx
chapters but are not found in docs/reference/. Use --strict to fail on any
difference.

Sections compared:
  - appendix/3 command examples  -> docs/reference/cli.md
  - appendix/2 2.1.1-2.1.3 keys  -> docs/reference/app-yaml.md
  - appendix/6 status codes      -> docs/reference/exit-codes.md
"""
import argparse
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
CN = REPO_ROOT / "docs" / "cn" / "source"
REF = REPO_ROOT / "docs" / "reference"

CLI_RE = re.compile(r"\bscalebox\s+([a-z][a-z0-9-]*)(?:\s+([a-z][a-z0-9-]*))?")
REF_CLI_RE = re.compile(r"^## `scalebox([^`]*)`")
PARAM_RE = re.compile(r"`([a-z][a-z0-9_]*)`")
TABLE_ROW_RE = re.compile(r"^\s*\|(.+)\|\s*$")
SECTION_RE = re.compile(r"^###\s+(2\.1\.\d|2\.1\.4|2\.2\.\d|2\.3)\b")

CLI_STOPWORDS = {
    "the", "and", "or", "to", "is", "in", "of", "for", "with",
    "command", "commands",
}


def read(path):
    return path.read_text(encoding="utf-8", errors="replace")


def parse_sphinx_cli(path):
    commands = set()
    for line in read(path).splitlines():
        for match in CLI_RE.finditer(line):
            top = match.group(1)
            sub = match.group(2)
            if top in CLI_STOPWORDS:
                continue
            if "xxxx" in top or "xxxx" in (sub or ""):
                continue
            parts = ["scalebox", top]
            if sub and not sub.startswith("-") and sub not in CLI_STOPWORDS:
                parts.append(sub)
            commands.add(" ".join(parts))
    return commands


def parse_ref_cli(path):
    commands = set()
    for line in read(path).splitlines():
        match = REF_CLI_RE.match(line)
        if match:
            parts = match.group(1).split()
            commands.add(" ".join(["scalebox"] + parts))
    return commands


def covered(command, reference):
    if command in reference:
        return True
    return any(item.startswith(command + " ") for item in reference)


def parse_appendix2_params(path):
    """Return the schema-uncovered parameters kept in appendix/2 §2.1.3."""
    params = set()
    section = ""
    for line in read(path).splitlines():
        match = SECTION_RE.match(line)
        if match:
            section = match.group(1)
            continue
        if section != "2.1.3":
            continue
        row = TABLE_ROW_RE.match(line)
        if row:
            cells = [cell.strip() for cell in row.group(1).split("|")]
            token = cells[0].strip().strip("`").strip()
            token = re.sub(r"^[-*]\s*", "", token)
            if re.fullmatch(r"[a-z][a-z0-9_]*", token):
                params.add(token)
    return params


def parse_ref_params(path):
    params = set()
    for line in read(path).splitlines():
        for token in PARAM_RE.findall(line):
            params.add(token)
        row = TABLE_ROW_RE.match(line)
        if row:
            cells = [cell.strip().strip("`") for cell in row.group(1).split("|")]
            token = cells[0]
            if re.fullmatch(r"[a-z][a-z0-9_]*", token):
                params.add(token)
    return params


def parse_codes(path):
    codes = set()
    for line in read(path).splitlines():
        row = TABLE_ROW_RE.match(line)
        if not row:
            continue
        cell = row.group(1).split("|")[0].strip().strip("`").strip()
        match = re.fullmatch(r"(\d{1,3})", cell)
        if match:
            codes.add(int(match.group(1)))
            continue
        match = re.fullmatch(r"(\d{1,3})\s*[–-]\s*(\d{1,3})", cell)
        if match:
            codes.update(range(int(match.group(1)), int(match.group(2)) + 1))
    return codes


def print_section(title, total, missing):
    print("== %s ==" % title)
    print("checked %d; missing %d" % (total, len(missing)))
    for item in missing:
        print("  %s" % item)
    print()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--strict", action="store_true", help="exit 1 when differences exist")
    args = parser.parse_args()

    sphinx_cli = parse_sphinx_cli(CN / "appendix" / "3_commandline_tools.md")
    ref_cli = parse_ref_cli(REF / "cli.md")
    missing_cli = sorted(cmd for cmd in sphinx_cli if not covered(cmd, ref_cli))

    sphinx_params = parse_appendix2_params(CN / "appendix" / "2_parameter_reference.md")
    ref_params = parse_ref_params(REF / "app-yaml.md")
    missing_params = sorted(sphinx_params - ref_params)

    sphinx_codes = parse_codes(CN / "appendix" / "6_exit_code_spec.md")
    ref_codes = parse_codes(REF / "exit-codes.md")
    missing_codes = sorted(sphinx_codes - ref_codes)

    print_section("CLI commands (appendix/3 vs reference/cli.md)", len(sphinx_cli), missing_cli)
    print_section("Schema-uncovered parameters (appendix/2 2.1.3 vs reference/app-yaml.md)",
                  len(sphinx_params), missing_params)
    print_section("Status codes (appendix/6 vs reference/exit-codes.md)",
                  len(sphinx_codes), missing_codes)

    if args.strict and (missing_cli or missing_params or missing_codes):
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
