import csv
from itertools import permutations
import re
import time
from collections import Counter

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
        count = row[1]
        if word in words:
            continue
        if not re.match('[a-z]+', word):
            continue
        if word not in valids:
            continue
        words[word] = count

# print(len(words))
# exit(1)

# we want a list of words that are all constructed of the same 6 letters

# words = ['asdfgh', 'asdf']
        
def issublist(check, base):
    countsCheck, countsBase = Counter(check), Counter(base)
    for letter, occs in countsCheck.items():
        if occs > countsBase[letter]:
            return False
    return True

def get_groups_func(full_words):
    groups = []
    i = 0
    start_time = time.perf_counter()
    for word in full_words:
        i += 1
        if i % 100 == 0:
            cur_time = time.perf_counter()
            print(i, 'time passed:', f'{cur_time-start_time:04f}', 'seconds. percent done: ', f'{i / len(full_words):04f}')
        group = [word]
        for lil_word in words:
            if issublist(lil_word, word):
                group.append(lil_word)
        groups.append(group)
    return groups

if __name__ == "__main__":

    full_words = list(filter(lambda w: len(w) == 6, words))
    print('full words:', len(full_words))

    groups = get_groups_func(full_words)

    with open('out_6.txt', 'w') as f:
        for g in groups:
            f.write(','.join(g))
            f.write('\n')