# pi-coding-agent LLM Setup Guide

**Last Updated:** March 11, 2026  
**pi Version:** 0.57.1

## ✅ Configuration Complete

Your pi-coding-agent is now configured with:
- **Default Provider:** Google Antigravity (free with Google account)
- **Default Model:** Gemini 3.1 Pro Low (1M context, minimal thinking = fast & free)
- **19+ Available Models:** Claude, GPT-5 Codex, GLM, Kimi, MiniMax, etc.
- **Optimized Settings:** Fast compaction, smart retry, minimal thinking by default

---

## 🔑 Current Authentication Status

### ✅ Working (OAuth)
| Provider | Status | Models Available |
|----------|--------|------------------|
| **Google Antigravity** | ✅ Active | Gemini 3.1 Pro, Claude Sonnet/Opus, GPT-OSS |
| **OpenAI Codex** | ✅ Active | GPT-5.x Codex models |
| **GitHub Copilot** | ✅ Active | Claude, GPT-5, Gemini |

### ⚠️ Needs API Key
| Provider | Key Needed | Get Key |
|----------|------------|---------|
| **MiniMax** | `MINIMAX_API_KEY` | https://platform.minimax.io/ |
| **Kimi (Moonshot)** | `KIMI_API_KEY` | https://platform.moonshot.ai/ |
| **Z.ai (GLM)** | `ZAI_API_KEY` | https://platform.z.ai/ |
| **OpenCode Zen** | `OPENCODE_API_KEY` | https://opencode.ai/ |

---

## 🚀 Quick Start

```bash
# Start pi with defaults (Gemini 3.1 Pro Low via Antigravity)
pi

# Test it works
pi --print "Hello!"

# Use specific model
pi --model claude-sonnet-4-6 --thinking low

# Use high-thinking mode for complex tasks
pi --model gemini-3.1-pro-high --thinking high

# List all available models
pi --list-models
```

---

## 📊 Available Models by Provider

### Google Antigravity (Free, OAuth)
| Model | Context | Thinking | Best For |
|-------|---------|----------|----------|
| `gemini-3.1-pro-low` | 1M | Yes | Default, fast tasks |
| `gemini-3.1-pro-high` | 1M | Yes | Complex reasoning |
| `gemini-3-flash` | 1M | Yes | Quick responses |
| `claude-sonnet-4-6` | 200K | Yes | Coding tasks |
| `claude-sonnet-4-5` | 200K | Yes | Daily coding |
| `claude-opus-4-6-thinking` | 200K | Yes | Complex architecture |
| `claude-opus-4-5-thinking` | 200K | Yes | Deep reasoning |

### GitHub Copilot (Subscription, OAuth)
| Model | Context | Thinking | Best For |
|-------|---------|----------|----------|
| `claude-sonnet-4.6` | 128K | Yes | Daily coding |
| `claude-opus-4.6` | 128K | Yes | Complex tasks |
| `gpt-5.4-codex` | 400K | Yes | Large codebases |
| `gpt-5.3-codex` | 272K | Yes | Code generation |

### OpenAI Codex (Subscription, OAuth)
| Model | Context | Thinking | Best For |
|-------|---------|----------|----------|
| `gpt-5.4-codex` | 272K | Yes | Latest Codex |
| `gpt-5.3-codex` | 272K | Yes | Code tasks |
| `gpt-5.2-codex` | 272K | Yes | Stable Codex |

### MiniMax (API Key Required)
| Model | Context | Thinking | Best For |
|-------|---------|----------|----------|
| `MiniMax-M2.5` | 205K | Yes | SOTA coding (80.2% SWE-Bench) |
| `MiniMax-M2.1` | 205K | Yes | Fast coding |
| `minimax-m2.5-free` | 205K | Yes | Free tier |

### Kimi/Moonshot (API Key Required)
| Model | Context | Thinking | Best For |
|-------|---------|----------|----------|
| `kimi-k2.5` | 262K | Yes | Agentic workflows |

### Z.ai/GLM (API Key Required)
| Model | Context | Thinking | Best For |
|-------|---------|----------|----------|
| `glm-5` | 205K | Yes | Complex systems |
| `glm-4.7` | 205K | Yes | Balanced tasks |
| `glm-4.6` | 205K | Yes | Quick tasks |

---

## 🔧 Adding API Keys

To use MiniMax, Kimi, or GLM models, add API keys to `~/.pi/agent/auth.json`:

```json
{
  "minimax": {
    "type": "api_key",
    "key": "your-minimax-api-key"
  },
  "kimi-coding": {
    "type": "api_key",
    "key": "your-kimi-api-key"
  },
  "zai": {
    "type": "api_key",
    "key": "your-zai-api-key"
  }
}
```

Then run:
```bash
chmod 600 ~/.pi/agent/auth.json
```

---

## ⚙️ Configuration Files

### `~/.pi/agent/settings.json`
Main configuration - models, thinking levels, behavior.

### `~/.pi/agent/auth.json`
Authentication credentials (permissions: 600).

### `~/.pi/agent/skills/`
Custom skills:
- `remotion-best-practices/` - Remotion video framework

---

## 💡 Usage Tips

### Thinking Levels
```bash
/thinking minimal    # 512 tokens - fastest (default)
/thinking low        # 2048 tokens
/thinking medium     # 8192 tokens
/thinking high       # 16384 tokens
/thinking xhigh      # 32768 tokens
```

### Model Switching (Interactive)
```
/model gemini-3.1-pro-low     # Fast default
/model claude-sonnet-4-6      # Best for coding
/model claude-opus-4-6-thinking  # Complex tasks
/model gpt-5.4-codex          # Large context
```

### Cost Optimization
1. Use `gemini-3.1-pro-low` with `minimal` thinking for quick tasks
2. Use `claude-sonnet-4-6` for daily coding
3. Reserve `claude-opus-4-6-thinking` for complex architecture
4. Compaction is enabled to reduce context tokens

---

## 🆘 Troubleshooting

### "No API key found for X"
Add the API key to `~/.pi/agent/auth.json` (see above).

### "Model not found"
Check exact model name with `pi --list-models`.

### "Authentication failed"
Re-authenticate with OAuth providers:
```bash
pi
/login google-antigravity
/login openai-codex
/login github-copilot
```

### Rate Limiting (Antigravity)
Google Antigravity has free tier rate limits. If hit:
1. Wait a few minutes
2. Switch to a different provider
3. Use paid API keys

---

## 📚 Resources

- [pi-coding-agent docs](/opt/pi-coding-agent/docs/)
- [Providers guide](/opt/pi-coding-agent/docs/providers.md)
- [Models guide](/opt/pi-coding-agent/docs/models.md)
- [Settings reference](/opt/pi-coding-agent/docs/settings.md)

### Provider Links
- [Google Antigravity](https://antigravity.dev/)
- [MiniMax Platform](https://platform.minimax.io/)
- [Moonshot Kimi](https://platform.moonshot.ai/)
- [Z.ai GLM](https://platform.z.ai/)

---

**Setup completed:** March 11, 2026  
**Tested working:** `pi --provider google-antigravity --model gemini-3.1-pro-low --print "Hello!"`
