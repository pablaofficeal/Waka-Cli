import subprocess
import datetime

def format_go_files():
    # находим только .go файлы
    go_files = subprocess.check_output(
        ['find', '.', '-name', '*.go', '-not', '-path', './vendor/*']
    ).decode().splitlines()

    for f in go_files:
        subprocess.call(['go', 'fmt', f])

subprocess.call(['sleep', '2'])

format_go_files()
status = subprocess.check_output(['git', 'status', '--porcelain']).decode().strip()
if not status:
    print("Нет изменений — коммит не нужен.")
    exit(0)

subprocess.check_call(['git', 'add', '.'])
msg = f'fmt: format go files {datetime.datetime.now().strftime("%Y-%m-%d %H:%M")}'
subprocess.check_call(['git', 'commit', '-m', msg])
subprocess.check_call(['git', 'push'])

print("✓ Коммит и пуш выполнены.")