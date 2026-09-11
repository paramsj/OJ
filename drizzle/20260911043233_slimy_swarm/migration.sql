CREATE TYPE "user_role" AS ENUM('USER', 'ADMIN');--> statement-breakpoint
CREATE TABLE "click_events" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	"short_link_id" uuid NOT NULL,
	"ip_address" text,
	"user_agent" text,
	"referrer" text,
	"clicked_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "id_ranges" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	"server_id" text NOT NULL UNIQUE,
	"start_id" bigint NOT NULL,
	"end_id" bigint NOT NULL,
	"next_id" bigint NOT NULL,
	"is_active" boolean DEFAULT true NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "refresh_tokens" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	"user_id" uuid NOT NULL,
	"token_hash" text NOT NULL,
	"expires_at" timestamp NOT NULL,
	"is_revoked" boolean DEFAULT false NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "short_links" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	"user_id" uuid NOT NULL,
	"original_url" text NOT NULL,
	"short_code" text NOT NULL UNIQUE,
	"title" text,
	"total_clicks" bigint DEFAULT 0 NOT NULL,
	"is_active" boolean DEFAULT true NOT NULL,
	"expires_at" timestamp,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "users" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	"name" text NOT NULL,
	"email" text NOT NULL UNIQUE,
	"password_hash" text NOT NULL,
	"role" "user_role" DEFAULT 'USER'::"user_role" NOT NULL,
	"is_active" boolean DEFAULT true NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE INDEX "idx_click_events_short_link_id" ON "click_events" ("short_link_id");--> statement-breakpoint
CREATE INDEX "idx_click_events_clicked_at" ON "click_events" ("clicked_at");--> statement-breakpoint
CREATE INDEX "idx_refresh_tokens_user_id" ON "refresh_tokens" ("user_id");--> statement-breakpoint
CREATE UNIQUE INDEX "idx_refresh_tokens_token_hash" ON "refresh_tokens" ("token_hash");--> statement-breakpoint
CREATE INDEX "idx_short_links_user_id" ON "short_links" ("user_id");--> statement-breakpoint
CREATE UNIQUE INDEX "idx_short_links_short_code" ON "short_links" ("short_code");--> statement-breakpoint
ALTER TABLE "click_events" ADD CONSTRAINT "click_events_short_link_id_short_links_id_fkey" FOREIGN KEY ("short_link_id") REFERENCES "short_links"("id") ON DELETE CASCADE;--> statement-breakpoint
ALTER TABLE "refresh_tokens" ADD CONSTRAINT "refresh_tokens_user_id_users_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;--> statement-breakpoint
ALTER TABLE "short_links" ADD CONSTRAINT "short_links_user_id_users_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;