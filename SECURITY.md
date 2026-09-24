# Security Policy / سياسة الأمان

## Supported version
Security fixes target the current `main` branch.

## Deployment model
HyperCache-Go has no authentication, TLS, ACLs, or persistence. Its default address is loopback (`127.0.0.1:6379`). Do not expose it directly to an untrusted network. Use network isolation and a secure transport/proxy when deployment requires remote access.

## Reporting
Please report suspected vulnerabilities privately through GitHub's security reporting facilities when available. Do not publish credentials, sensitive cached values, or exploit details in a public issue.

## العربية
يدعم الفرع `main` الحالي إصلاحات الأمان. لا يحتوي المشروع حاليًا على مصادقة أو TLS أو ACL، لذلك لا ينبغي كشفه مباشرة لشبكة غير موثوقة. استخدم قناة الإبلاغ الأمني الخاصة في GitHub عند توفرها، ولا تنشر أسرارًا أو بيانات حساسة في Issue عامة.

Maintainer: Radwan Abdulhadi Ahmed / رضوان عبدالهادي أحمد / @rad03i2
