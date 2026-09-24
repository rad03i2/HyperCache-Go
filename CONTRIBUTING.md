# Contributing / المساهمة

Thanks for improving HyperCache-Go. Keep changes focused, dependency-light, and covered by tests.

1. Create a branch from `main`.
2. Run `go vet ./...`, `go test -race ./...`, and `go build ./cmd/hypercache`.
3. Document user-visible behavior and protocol changes.
4. Never commit credentials, production cache data, generated binaries, or personal information.
5. Open a focused pull request describing behavior and validation.

## العربية

نرحب بالمساهمات المركزة والقابلة للاختبار. شغّل فحوصات `go vet` والاختبارات مع race detector والبناء قبل إرسال Pull Request، ووثّق أي تغيير ظاهر للمستخدم أو للبروتوكول. لا ترفع أسرارًا أو بيانات كاش حقيقية أو ملفات تنفيذية مولدة أو معلومات شخصية.

Maintainer: Radwan Abdulhadi Ahmed / رضوان عبدالهادي أحمد / @rad03i2
