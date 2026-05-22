import * as fs from 'fs';
import * as path from 'path';
import Ajv, { ErrorObject } from 'ajv';
import addFormats from 'ajv-formats';

/**
 * Validates a payload against a platform-specific JSON schema.
 * Uses ajv for full JSON Schema draft-07 support.
 */
export class SchemaValidator {
  private schemaDir: string;
  private ajv: Ajv;
  private validators: Map<string, any>;

  constructor() {
    this.schemaDir = path.resolve(__dirname, '../schemas');
    this.ajv = new Ajv({
      allErrors: true,
      strict: false,         // 宽松模式，允许 schema 中有 ajv 不认识的 keyword
      verbose: true,          // 给出更详细的错误上下文
    });
    addFormats(this.ajv);   // 支持 format: email/uri/date 等
    this.validators = new Map();
  }

  /**
   * Validate payload against platform schema.
   * Falls back to basic checks if schema missing or ajv compile fails.
   */
  public async validate(platform: string, type: string, data: any): Promise<{ valid: boolean; errors?: string[] }> {
    const platformKey = platform.toLowerCase();
    const typeKey = type.replace(/-([a-z])/g, (_: any, c: string) => c.toUpperCase());
    const schemaFile = path.join(this.schemaDir, 'platforms', `${platformKey}.${typeKey}.schema.json`);

    // ── Case 1: No schema file → basic check + warning ─────────────
    if (!fs.existsSync(schemaFile)) {
      console.warn(
        `[Validator] ⚠️  No schema found for ${platform} (${type}) at ${schemaFile}.\n` +
        `             Falling back to basic checks. ` +
        `Consider adding a schema to improve validation.`
      );
      return this.basicValidate(data);
    }

    // ── Case 2: Schema exists → full ajv validation ────────────────
    try {
      let validate = this.validators.get(schemaFile);
      if (!validate) {
        const schema = JSON.parse(fs.readFileSync(schemaFile, 'utf-8').replace(/^\uFEFF/, ''));

        // Remove unsupported keywords gracefully (ajv may not support all custom keywords)
        this.sanitizeSchema(schema);

        validate = this.ajv.compile(schema);
        this.validators.set(schemaFile, validate);
      }

      // ── Auto-detect validation target ──────────────────────────────
      // Schema validates `contentPublishForm` (inner object).
      // Payload may be:
      //   A) { accountForms: [ { platformAccountId, contentPublishForm: {...} ] }
      //   B) { formType, title, description, ... }  (already inner object)
      // We detect and validate the inner object(s) accordingly.
      let targets: any[] = [];

      if (data.accountForms && Array.isArray(data.accountForms)) {
        // Case A: full payload wrapper → validate each contentPublishForm
        targets = data.accountForms
          .map((f: any) => f.contentPublishForm)
          .filter((c: any) => c);
      } else if (data.formType || data.title || data.description) {
        // Case B: already inner object
        targets = [data];
      } else {
        // Fallback: validate the whole data as-is
        targets = [data];
      }

      let allErrors: string[] = [];
      for (let i = 0; i < targets.length; i++) {
        const target = targets[i];
        const idx = i;
        const valid = validate(target);
        if (!valid && validate.errors) {
          const prefix = targets.length > 1 ? `accountForms[${idx}].contentPublishForm: ` : '';
          const errs = this.formatAjvErrors(validate.errors).map(e => prefix + e);
          allErrors.push(...errs);
        }
      }

      if (allErrors.length > 0) {
        return { valid: false, errors: allErrors };
      }

      return { valid: true };
    } catch (err: any) {
      console.warn(`[Validator] ⚠️  ajv compile/validate failed: ${err.message}`);
      console.warn(`             Falling back to basic checks.`);
      return this.basicValidate(data);
    }
  }

  // ── Full ajv validation error formatter ─────────────────────────────
  private formatAjvErrors(errors: ErrorObject[]): string[] {
    return errors.map((e) => {
      const instancePath = e.instancePath || '/';
      const message = e.message || 'validation error';
      const params = e.params as any;

      // Nicer messages for common errors
      if (e.keyword === 'enum') {
        const allowed = (params?.allowedValues || []).join(', ');
        return `${instancePath}: ${message} (allowed: ${allowed})`;
      }
      if (e.keyword === 'pattern') {
        return `${instancePath}: ${message} (pattern: ${params?.pattern || '?'})`;
      }
      if (e.keyword === 'required') {
        const missing = (params?.missingProperty) ? `${params.missingProperty}` : '';
        return `${instancePath}: missing required field "${missing}"`;
      }
      if (e.keyword === 'additionalProperties') {
        const extra = params?.additionalProperty || '?';
        return `${instancePath}: unexpected field "${extra}" (not in schema)`;
      }

      return `${instancePath}: ${message}`;
    });
  }

  // ── Remove keywords ajv doesn't understand ──────────────────────────
  private sanitizeSchema(schema: any): void {
    const unsupported = ['elys'];// add more if needed
    if (typeof schema === 'object' && schema !== null) {
      for (const key of unsupported) {
        delete schema[key];
      }
      for (const value of Object.values(schema)) {
        if (typeof value === 'object' && value !== null) {
          this.sanitizeSchema(value);
        }
      }
    }
  }

  // ── Basic fallback validation (when schema missing or ajv fails) ───
  private basicValidate(data: any): { valid: boolean; errors?: string[] } {
    const errors: string[] = [];

    // 1. Check accountForms exists and is an array
    if (data.formType || data.title || data.description || data.content) {
      if (data.content === undefined && data.description === undefined) {
        errors.push('Inner form payload must include content or description');
      }
    } else if (!data.accountForms) {
      errors.push('Missing required field: accountForms');
    } else if (!Array.isArray(data.accountForms)) {
      errors.push('accountForms must be an array');
    } else {
      // 2. Per-accountForm checks
      data.accountForms.forEach((form: any, idx: number) => {
        if (!form.platformAccountId && !form.account_id) {
          errors.push(`accountForms[${idx}]: missing platformAccountId`);
        }
        if (form.contentPublishForm) {
          const cpf = form.contentPublishForm;
          // Check images for image-text
          if (data.publishType === 'image-text' || !data.publishType) {
            if (!form.images || !Array.isArray(form.images) || form.images.length === 0) {
              errors.push(`accountForms[${idx}]: images array is required for image-text`);
            }
          }
          // Check video for video
          if (data.publishType === 'video') {
            if (!form.video) {
              errors.push(`accountForms[${idx}]: video object is required for video publish`);
            }
          }
          // Check content for article
          if (data.publishType === 'article') {
            if (!cpf.content) {
              errors.push(`accountForms[${idx}].contentPublishForm: content is required for article publish`);
            }
          }
        } else {
          errors.push(`accountForms[${idx}]: missing contentPublishForm`);
        }
      });
    }

    return {
      valid: errors.length === 0,
      errors: errors.length > 0 ? errors : undefined,
    };
  }
}
