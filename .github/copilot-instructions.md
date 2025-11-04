# Copilot Instructions for terraform-provider-hyperv
- **Big picture**: This Terraform provider drives Hyper-V by rendering PowerShell templates and executing them remotely via WinRM or optional SSH, keeping Terraform state in sync with the host.
- **Code layout**: `internal/provider` contains Terraform schemas/resources, `api/` houses Hyper-V domain models plus remote client abstractions, and `powershell/` packages utility templates for command execution.
- **Data flow**: Resource CRUD functions type assert `meta.(api.Client)` and call helpers in `api/hyperv` (e.g. `CreateVm`), which render text/templates and rely on `RunScriptWithResult` to unmarshal JSON back into Go structs.
- **Remote execution contract**: Every PowerShell template in `api/hyperv` must finish by emitting `ConvertTo-Json` output; both `winrm_helper` and `ssh_helper` unmarshal that string directly.
- **Connection management**: WinRM sessions are pooled in `internal/provider/config.go` with `go-commons-pool` (MaxTotal=5); always borrow/return via the helper methods instead of instantiating raw clients.
- **SSH mode**: Setting the provider `ssh` flag routes through `api/ssh-helper`, which still executes PowerShell on Windows hosts (`IsWindows=true`); keep templates Windows-friendly even when testing via SSH.
- **Resource pattern**: Follow the lifecycle structure in `internal/provider/resource_hyperv_machine_instance.go`—create/update call the client, then re-read state, with polling helpers like `turnOffVmIfOn` and `waitForIps` handling long operations.
- **Validation helpers**: Reuse functions in `internal/provider/schema_validators.go` (`AllowedIsoVolumeName`, `StringKeyInMap`, etc.) so diagnostics stay consistent and localized.
- **Enumerations**: Map Terraform strings through the enum tables in `api/*.go` (e.g. `api.StartAction_name`) before sending numeric values to PowerShell scripts.
- **File transfer**: Rely on the abstractions in `api.Client` (`UploadFile`, `UploadDirectory`, `DeleteFileOrDirectory`) instead of shelling out—see `resource_hyperv_iso_image.go` for examples of ensuring remote file state.
- **Diff suppression**: Normalize Windows paths the way `resource_hyperv_machine_instance.go` does (replace `\` with `/`, compare case-insensitively) whenever detecting drift on host paths.
- **Unit tests**: Run `mise run test` (alias for `go test ./... -short -timeout 10m`) to cover JSON marshaling/parsing in `api/` when touching templates or models.
- **Acceptance tests**: Use `mise run testacc` or `TF_ACC=1 go test ./... -v -timeout 120m`; they require a reachable Hyper-V host with valid WinRM/SSH credentials provided via `HYPERV_*` environment variables.
- **Bootstrap**: Execute `mise run setup` once per clone to download modules and generate a `.terraformrc` that points Terraform at the locally built provider.
- **Build/install**: `mise run install` (or `make install`) compiles the binary and drops it in `GOBIN`; Terraform reads its location from the generated `.terraformrc`.
- **Dev environment**: `mise run dev-init` and `mise run dev-plan` operate on the configs in `dev/`, assuming the provider is installed and `.terraformrc` is current.
- **Formatting**: Run `mise run fmt` to execute both `go fmt ./...` and `terraform fmt` on the examples, and `mise run lint` to invoke `golangci-lint` using `.golangci.yml`.
- **Documentation**: After schema changes, run `mise run docs` to regenerate content under `docs/` via `tfplugindocs`; the directives live in `main.go`.
- **Logging**: The provider removes default timestamps (`log.SetFlags` in `main.go`); rely on the `log.Printf` calls in `winrm_helper` and `ssh_helper` for tracing remote execution.
- **Debug mode**: Launch `go run . -debug` to expose the provider to Terraform debuggers (set `TF_REATTACH_PROVIDERS` per Terraform docs) when troubleshooting complex flows.
- **Connection settings**: When adding schema attributes for auth, update both `provider.go` definitions and the `Config` struct wiring so WinRM, SSH, and `.terraformrc.tpl` stay aligned.
- **Concurrency/timeouts**: Honor the timeout constants declared at the top of each resource file; long-running operations should respect Terraform `ResourceTimeout` plus the helper polling loops.
- **Adding API calls**: Extend `api.Client` in `api/provider.go`, implement the method on `hyperv.ClientConfig`, and provide matching logic in both client helpers; new features need PowerShell templates under `api/hyperv/`.
- **PowerShell templates**: Define templates with `template.Must`, sanitize input (see `powershell/template.go` helpers), and ensure they clean up temp files when errors occur.
- **Elevation**: PowerShell commands can be wrapped as scheduled tasks via `powershell.RunPowershell`; design scripts to tolerate that execution context and remove temporary artifacts.
- **File handling**: Large transfers convert files to base64 chunks (`powershell/template.go`); avoid rewriting this mechanism to prevent memory spikes on the host.
- **Examples/docs**: Mirror resource attributes across `docs/` and `examples/` directories when making user-visible changes so the Terraform Registry docs stay accurate.
- **Release metadata**: The provider address is pinned to `registry.terraform.io/bafbi/hyperv` (see `main.go` and `terraform-registry-manifest.json`); coordinate any change with registry updates.
- **CI expectations**: GitHub Actions assume `mise` commands exist; add new workflows by updating `mise.toml` so local dev and CI stay in sync.

## Sandbox Hyper-V host (testing)

- A reusable sandbox host is available for quick PowerShell/WinRM testing: `neudeline@172.16.0.220` (SSH).
- I verified connectivity and PowerShell execution non-interactively from the development machine. Example used:

```sh
ssh neudeline@172.16.0.220 powershell -NoProfile -Command "Write-Output 'hello-from-sandbox'"
```

Output from the verification was: `hello-from-sandbox`.

- Notes and safety:
	- The sandbox runs a Windows host and accepts SSH connections; the provider's SSH mode (`ssh=true`) still executes PowerShell on the remote host (`IsWindows=true`).
	- Prefer using SSH key auth (set `ssh_private_key_path` or `ssh_private_key`) and avoid checking credentials into repo files. The repo contains an `.terraformrc.tpl` that references `GOBIN` if you need to wire the locally-built provider into Terraform dev runs.
	- For scripted tests, use non-destructive PowerShell commands and always clean up resources created during acceptance testing (or run with `TF_ACC=1` and a dedicated sandbox account).

Add or remove this sandbox entry if the host changes or access is revoked.
