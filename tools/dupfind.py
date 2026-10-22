#!/usr/bin/env python3
# Reference implementation of the DupFind contract, used only to validate the grader.
import csv, hashlib, os, sys
def usage(m): print(m, file=sys.stderr); sys.exit(2)
args = sys.argv[1:]; summary = False; min_size = 1
while args and args[0].startswith('-'):
    a = args.pop(0)
    if a == '--summary': summary = True
    elif a == '--min-size':
        if not args: usage('missing value')
        v = args.pop(0)
        if not v.isdigit() or int(v) < 1: usage('bad min-size')
        min_size = int(v)
    else: usage('unknown option ' + a)
if len(args) != 1: usage('need one directory')
root = args[0]
if not os.path.isdir(root): print('not a directory', file=sys.stderr); sys.exit(1)
by_size = {}
for d, dirs, files in os.walk(root):
    for n in files:
        p = os.path.join(d, n)
        if os.path.islink(p) or not os.path.isfile(p): continue
        s = os.path.getsize(p)
        if s < min_size: continue
        by_size.setdefault(s, []).append(os.path.relpath(p, root).replace(os.sep, '/'))
groups = []
for s, ps in by_size.items():
    if len(ps) < 2: continue
    by_hash = {}
    for p in ps:
        h = hashlib.sha256(open(os.path.join(root, p), 'rb').read()).hexdigest()
        by_hash.setdefault(h, []).append(p)
    for h, hp in by_hash.items():
        if len(hp) > 1: groups.append((s, h, sorted(hp, key=lambda x: x.encode())))
groups.sort(key=lambda g: (-g[0] * (len(g[2]) - 1), g[1]))
w = csv.writer(sys.stdout, lineterminator='\n')
if summary:
    w.writerow(['groups', 'duplicate_files', 'reclaimable_bytes'])
    w.writerow([len(groups), sum(len(g[2]) - 1 for g in groups), sum(g[0] * (len(g[2]) - 1) for g in groups)])
else:
    w.writerow(['group', 'size', 'sha256', 'path'])
    for i, (s, h, ps) in enumerate(groups, 1):
        for p in ps: w.writerow([i, s, h, p])
