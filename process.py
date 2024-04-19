import csv
from collections import Counter
import re

valids = set()
with open("valids.csv") as fp:
    reader = csv.reader(fp, delimiter=",")
    for row in reader:
        valids.add(row[0].lower())

words = {}
with open("words.csv") as fp:
    reader = csv.reader(fp, delimiter=",", quotechar='"')
    for row in reader:
        word = row[0]
        count = int(row[1])
        if word in words:
            continue
        if not re.match('[a-z]+', word):
            continue
        if word not in valids:
            continue
        words[word] = count

with open('out_6.txt') as fp:
    reader = csv.reader(fp)
    data = [row for row in reader]

data.sort(key = lambda l: len(l), reverse=True)
data.sort(key = lambda l: sum([words[w] for w in l]), reverse=True)

levels = []
for g in data:
    if len(g) < 30:
        continue
    base = g[0]
    g.sort(key=lambda x: words[x], reverse=True)
    levels.append(g[:20])
    if base not in levels[:-1]:
        levels[-1][-1] = base

print(len(levels))
for level in levels[:100]:
    print(level)

# print(data[:10])

# counted = Counter([len(l) for l in data])
# for c in counted.items():
#     print(c)