#!/usr/bin/env python3
"""Replace template.ParseFiles blocks with renderTemplate calls."""
import re

files_config = [
    ('go/login.go', [
        ('html/library.html',  False),
        ('html/about.html',    False),
        ('html/login.html',    False),
        ('html/admin.html',    'librarysum'),
        ('html/ranking.html',  False),
    ]),
    ('go/adminbooks.go', [
        ('html/addbook.html',           False),
        ('html/view-book.html',         False),
        ('html/view-lend-records.html',  False),
        ('html/view-return-records.html', False),
        ('html/adjust-book.html',        False),
    ]),
    ('go/adminnotices.go', [
        ('html/addnotice.html',   False),
        ('html/view-notice.html',  False),
    ]),
    ('go/userbook.go', [
        ('html/lend-book-list.html',    False),
        ('html/view-adjustbook.html',   False),
    ]),
    ('go/users.go', [
        ('html/view-user.html',    False),
        ('html/view-useropi.html',  False),
        ('html/user-library.html',  False),
    ]),
]


def parsefiles_line(tmpl):
    return 'tmpl, err := template.ParseFiles("' + tmpl + '")'


def find_outer_if(content, pos):
    """Walk backward from pos to find the containing 'if' statement line start."""
    line_start = content.rfind('\n', 0, pos) + 1
    rpos = line_start - 2
    while rpos >= 0:
        if content[rpos:rpos+3] == 'if ':
            if_start = content.rfind('\n', 0, rpos) + 1
            if_line = content[if_start:content.find('\n', rpos)].strip()
            # Must end with { or contain ) {
            if if_line.endswith('{'):
                return if_start
        rpos -= 1
    return line_start


def find_block_end(content, brace_start):
    """Find matching closing brace."""
    depth = 1
    p = brace_start + 1
    while p < len(content) and depth > 0:
        if content[p] == '{':
            depth += 1
        elif content[p] == '}':
            depth -= 1
        p += 1
    # Check for 'return' after the block
    after = content[p:p+15].strip()
    if after.startswith('return'):
        p = content.find('\n', p) + 1
    return p


def run():
    for filepath, templates in files_config:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        original = content

        for tmpl_name, data in templates:
            search = parsefiles_line(tmpl_name)
            idx = 0
            while True:
                pos = content.find(search, idx)
                if pos < 0:
                    break

                if_start = find_outer_if(content, pos)
                brace_start = content.index('{', if_start)
                block_end = find_block_end(content, brace_start)

                indent = content[if_start:if_start + len(content[if_start:]) - len(content[if_start:].lstrip())]
                if data:
                    repl = f'{indent}renderTemplate(w, "{tmpl_name}", {data})\n'
                else:
                    repl = f'{indent}renderTemplate(w, "{tmpl_name}", nil)\n'

                content = content[:if_start] + repl + content[block_end:]
                idx = if_start + len(repl)

        if content != original:
            # Remove unused imports
            has_fmt = 'fmt.' in content or 'fmt.Printf' in content
            has_template = 'template.' in content
            lines = content.split('\n')
            new_lines = []
            for line in lines:
                s = line.strip().strip('"')
                if s == 'fmt' and not has_fmt:
                    continue
                if s == 'text/template' and not has_template:
                    continue
                new_lines.append(line)
            content = '\n'.join(new_lines)

            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f'  Updated {filepath}')
        else:
            print(f'  No changes in {filepath}')


if __name__ == '__main__':
    run()
