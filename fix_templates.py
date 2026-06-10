#!/usr/bin/env python3
"""Replace template.ParseFiles blocks with renderTemplate calls in Go handler files."""
import re
import sys

files_to_fix = {
    'go/login.go': [
        ('html/library.html', None),
        ('html/about.html', None),
        ('html/login.html', None),
        ('html/admin.html', 'librarysum'),
        ('html/ranking.html', None),
    ],
    'go/adminbooks.go': [
        ('html/addbook.html', None),
        ('html/view-book.html', None),
        ('html/view-lend-records.html', None),
        ('html/view-return-records.html', None),
        ('html/adjust-book.html', None),
    ],
    'go/adminnotices.go': [
        ('html/addnotice.html', None),
        ('html/view-notice.html', None),
    ],
    'go/userbook.go': [
        ('html/lend-book-list.html', None),
        ('html/lend-book-list.html', None),  # handleHTMLRequest
        ('html/view-adjustbook.html', None),
    ],
    'go/users.go': [
        ('html/view-user.html', None),
        ('html/view-useropi.html', None),
        ('html/user-library.html', None),
    ],
}

for filepath, templates in files_to_fix.items():
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    original = content

    for template_name, data_var in templates:
        # Pattern: tmpl, err := template.ParseFiles("html/xxx.html") ... error handling ... ExecuteTemplate ... error handling ... closing }
        # We need to find the exact block for each template

        # Find tmpl, err := template.ParseFiles("template_name")
        search = f'tmpl, err := template.ParseFiles("{template_name}")'
        idx = content.find(search)
        if idx < 0:
            print(f"  Skipping {template_name} (not found in {filepath})")
            continue

        # Find the start of the line
        line_start = content.rfind('\n', 0, idx) + 1
        indent = content[line_start:idx]

        # The replacement depends on whether data is provided
        if data_var:
            replacement = f'{indent}renderTemplate(w, "{template_name}", {data_var})'
        else:
            replacement = f'{indent}renderTemplate(w, "{template_name}", nil)'

        # Find the end of this if block by brace counting
        block_start = content.find('{', idx) + 1
        depth = 1
        pos = block_start
        while pos < len(content) and depth > 0:
            if content[pos] == '{':
                depth += 1
            elif content[pos] == '}':
                depth -= 1
            pos += 1

        block_end = pos  # position after the closing }

        # The full region to replace starts from the template.ParseFiles line
        # and ends at the closing }

        # Now we need to figure out what to replace:
        # Usually the pattern is:
        #   if r.Method == http.MethodGet {
        #       tmpl, err := template.ParseFiles(...)
        #       if err != nil { ... }
        #       err = tmpl.ExecuteTemplate(...)
        #       if err != nil { ... }
        #   }
        # We want to replace from the 'tmpl' line to the closing '}'

        # Find the end of the if block containing this template call
        # Walk backward from idx to find the if, then find the matching closing brace
        if_line_start = content.rfind('\n', 0, line_start - 2) + 1
        if_line = content[if_line_start:line_start]

        if 'if r.Method == http.MethodGet' in if_line or 'if role == "Admin"' in if_line:
            # This is the GET if block - we can replace the whole thing
            # Find the matching closing brace
            brace_pos = content.find('{', if_line_start)
            depth = 1
            pos = brace_pos + 1
            while pos < len(content) and depth > 0:
                if content[pos] == '{':
                    depth += 1
                elif content[pos] == '}':
                    depth -= 1
                pos += 1

            # The block ends at pos
            block_text = content[if_line_start:pos]

            # Check what follows - sometimes there's a return statement
            rest = content[pos:].lstrip()
            if rest.startswith('return'):
                # Eat the return statement too
                return_end = pos + content[pos:].find('\n', content[pos:].find('return'))
                pos = return_end + 1

            # Now we need to handle the indentation of the closing brace
            # The closing brace is at pos-1 in the content
            new_block = f'{indent}renderTemplate(w, "{template_name}", nil)\n'

            # Check if this if block is inside another check (like role == "Admin")
            context_before = content[max(0, if_line_start-100):if_line_start]
            if 'if role == "Admin"' in context_before or 'if role' in context_before:
                # Keep the if role check but replace the inner if method get
                inner_indent = indent + '\t'
                new_block = f'{inner_indent}renderTemplate(w, "{template_name}", nil)\n'

            content = content[:if_line_start] + new_block + content[pos:]
        else:
            # Just replace the template.ParseFiles block inline
            new_block = f'{indent}renderTemplate(w, "{template_name}", nil)\n'
            content = content[:line_start] + new_block + content[block_end:]

        print(f"  Replaced {template_name} in {filepath}")

    if content != original:
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"  ✓ {filepath} updated")
    else:
        print(f"  - {filepath} unchanged")

print("\nDone!")
