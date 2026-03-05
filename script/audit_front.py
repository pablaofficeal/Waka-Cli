import subprocess
import os
import time
import zipfile

projects = ["main/frontend", "admin/frontend", "main/plugin", "partner/frontend"]

AUDIT_DIR = "audit"
os.makedirs(AUDIT_DIR, exist_ok=True)

json_reports = []  # список всех путей к JSON-файлам

for project_path in projects:
    print(f"\n=== Обработка: {project_path} ===")
    
    if not os.path.exists(project_path):
        print(f"Ошибка: Директория {project_path} не найдена.")
        continue

    result = subprocess.run(
        ["npm", "audit", "--json"],
        cwd=project_path,
        capture_output=True,
        text=True
    )

    timestamp = time.strftime('%Y%m%d-%H%M%S')
    json_filename = f"npm_audit_{timestamp}_{project_path.replace('/', '_')}.json"
    json_path = os.path.join(AUDIT_DIR, json_filename)

    # сохраняем JSON независимо от фикса → это важно!
    with open(json_path, "w") as f:
        f.write(result.stdout)

    json_reports.append(json_path)

    if result.returncode != 0:
        print(f"[!] Найдены уязвимости в {project_path}. Пытаемся исправить...")

        subprocess.run(["npm", "audit"], cwd=project_path)

        print(f"\n[*] Ожидание 2 секунды перед npm audit fix...")
        time.sleep(2)

        subprocess.run(["npm", "audit", "fix"], cwd=project_path)
        subprocess.run(["npm", "audit", "fix", "--force"], cwd=project_path)

    else:
        print(f"[OK] Уязвимостей не обнаружено.")

# ============================================================
# СОЗДАЁМ ОДИН ZIP АРХИВ ИЗ ВСЕХ JSON ФАЙЛОВ
# ============================================================

zip_name = f"audit_zip_{time.strftime('%Y%m%d-%H%M%S')}.zip"
zip_path = os.path.join(AUDIT_DIR, zip_name)

with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as zipf:
    for report in json_reports:
        zipf.write(report, arcname=os.path.basename(report))

print(f"\n🎉 Все JSON-файлы упакованы в архив: {zip_path}")
print("Все задачи выполнены!")