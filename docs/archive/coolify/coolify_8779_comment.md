Thank you for maintaining such a high-quality project!

My PR #8779 was closed by the quality gate with the `quality/rejected` label, but I'd like to understand the specific failures so I can address them properly.

### Request for Clarification:

1. **Which specific quality checks failed?** The automated message doesn't include detailed logs.

2. **Are CI logs accessible?** I'd like to review the exact test failures or linting errors.

3. **What are the quality requirements** for SSH-related fixes in Coolify?

### Context:

This PR addresses the intermittent "Permission denied (publickey,password)" error (issue #7724) by ensuring proper SSH key file permissions (600) before use.

### Commitment:

I want to resubmit with a fix that **passes all quality gates**. Any guidance on:
- Required test coverage
- Documentation expectations
- Code style/formatting requirements (Prettier, ESLint, etc.)

would be greatly appreciated!

Thanks for your time and for building such an excellent deployment platform! 🙏

---

**Willing to:** Add tests, update docs, fix formatting, or provide additional context as needed.
