import os
import glob

base_path = r'e:\_src\golang\go-knowledge-base'

def get_primary_files(category, folder, count):
    files = []
    for i in range(1, count + 1):
        num = f'{i:03d}'
        pattern = os.path.join(base_path, folder, f'{category}-{num}-*.md')
        matches = glob.glob(pattern)
        if matches:
            # Get the largest file as primary
            primary = max(matches, key=os.path.getsize)
            size = os.path.getsize(primary)
            files.append({
                'cat': category,
                'num': num,
                'name': os.path.basename(primary),
                'size_kb': round(size / 1024, 1),
                'path': primary
            })
    return files

# Get all files
ec_files = get_primary_files('EC', '03-Engineering-CloudNative', 20)
ld_files = get_primary_files('LD', '02-Language-Design', 15)
ts_files = get_primary_files('TS', '04-Technology-Stack', 15)

all_files = ec_files + ld_files + ts_files

print('Category | Num | Size(KB) | Filename')
print('-' * 80)
for f in all_files:
    name = f['name'][:50]
    print(f"{f['cat']:8} | {f['num']} | {f['size_kb']:8} | {name}")

print(f'\nTotal files: {len(all_files)}')

# Save file list for later
with open('target_files.txt', 'w') as out:
    for f in all_files:
        out.write(f"{f['cat']}|{f['num']}|{f['path']}|{f['name']}\n")
