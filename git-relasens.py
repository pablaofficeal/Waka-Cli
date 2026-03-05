import subprocess


print("Starting release process...")

subprocess.run(["git", "checkout", "Release"])

subprocess.run(["sleep", "2"])

print("Merging main into Release...")

subprocess.run(["sleep", "2"])

subprocess.run(["git", "merge", "main"])

subprocess.run(["sleep", "2"])

print("Pushing Release...")

subprocess.run(["git", "push"])

print("Checking out main...")

subprocess.run(["sleep", "2"])

subprocess.run(["git", "checkout", "main"])

print("Release process completed.")