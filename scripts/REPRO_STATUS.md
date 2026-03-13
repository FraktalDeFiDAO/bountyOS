# 🔍 TLSX ARM64 Deadlock Reproduction Status

**Script:** `scripts/repro_arm64_deadlock.sh`  
**Purpose:** Forensic reproduction of the ARM64 deadlock issue in tlsx  
**Date:** March 13, 2026

---

## ✅ SCRIPT STATUS

**Fixed Issue:**
- ❌ **Bug:** `'quiet'` kernel parameter was being treated as a file
- ✅ **Fix:** Removed `'quiet'` from kernel append line

**Current Status:** Running in QEMU ARM64 emulation

---

## 📋 WHAT THE SCRIPT DOES

### Phase 1: Source Acquisition ✅
```bash
git clone --depth 1 https://github.com/projectdiscovery/tlsx temp_repro_tlsx
```
Clones the latest tlsx source code.

### Phase 2: Exploit Injection ✅
```bash
cp au-workspace/projects/bounty-tlsx-hangs/source/internal/pdcp/race_test.go temp_repro_tlsx/internal/pdcp/
```
Injects the PoC test that triggers the race condition.

### Phase 3: Cross-Compilation ✅
```bash
GOARCH=arm64 go test -c ./internal/pdcp -o tools/vm/pdcp_arm64.test
```
Cross-compiles the test for ARM64 architecture.

### Phase 4: Initrd Packaging ✅
```bash
docker run --rm --platform linux/amd64 -v "$(pwd):/work" ...
```
Packages the binary into an AArch64 boot filesystem (initramfs).

### Phase 5: AArch64 System Emulation 🔄
```bash
docker run --rm -v "$(pwd):/vm" -w /vm docker.io/tianon/qemu \
  qemu-system-aarch64 -machine virt -cpu cortex-a57 -m 512M \
  -kernel tools/vm/vmlinuz -initrd tools/vm/initramfs_repro.gz \
  -nographic -append 'console=ttyAMA0'
```
Boots a full ARM64 Linux VM using QEMU to run the test.

---

## ⏱️ EXPECTED RUNTIME

| Phase | Duration | Status |
|-------|----------|--------|
| Phase 1: Clone | ~30 seconds | ✅ Complete |
| Phase 2: Inject | ~5 seconds | ✅ Complete |
| Phase 3: Compile | ~60 seconds | ✅ Complete |
| Phase 4: Package | ~30 seconds | ✅ Complete |
| Phase 5: Emulate | **5-15 minutes** | 🔄 Running |

**Total Expected Time:** 10-20 minutes

---

## 🎯 WHAT TO EXPECT

When the QEMU VM boots, you should see:

1. **Linux boot messages** (ARM64 kernel)
2. **Init system startup**
3. **Test execution:**
   ```
   --- BOOTING ARM64 FORENSIC ENVIRONMENT ---
   === RUN   TestUploadWriterExploit
   ```
4. **Race condition output** (if reproduced)
5. **VM shutdown** (`poweroff -f`)

---

## 🛑 TO STOP THE SCRIPT

```bash
# In the terminal where it's running
Ctrl+C

# Or in another terminal
docker ps | grep qemu
docker stop <container-id>
```

---

## 📊 OUTPUT FILES

After completion, these files are created:

| File | Location | Purpose |
|------|----------|---------|
| Test Binary | `tools/vm/pdcp_arm64.test` | ARM64 compiled test |
| Initramfs | `tools/vm/initramfs_repro.gz` | Boot filesystem |
| Source Clone | `temp_repro_tlsx/` | Fresh tlsx source |

---

## 🔧 ALTERNATIVE: FASTER REPRODUCTION

If QEMU emulation is too slow, you can test directly:

```bash
# On an ARM64 machine (Raspberry Pi, M1/M2 Mac, etc.)
cd temp_repro_tlsx
GOARCH=arm64 go test -race -run TestUploadWriterExploit ./internal/pdcp -v

# Or use Docker with ARM64 emulation
docker run --rm --platform linux/arm64 -v $(pwd):/app -w /app golang:1.21 \
  go test -race -run TestUploadWriterExploit ./internal/pdcp -v
```

---

## 📝 NOTES

- **Why ARM64?** The original deadlock was reported on ARM64 systems
- **Why QEMU?** Provides architectural fidelity without physical hardware
- **Why Alpine?** Minimal initramfs for fast boot and clear output
- **Why 512MB RAM?** Minimum for reliable reproduction

---

**Status:** Script is running. Expected completion in 5-15 minutes.
