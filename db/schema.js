import { pgTable, uuid, text, timestamp, boolean, bigint , pgEnum , index, uniqueIndex } from "drizzle-orm/pg-core";

export const userRoleEnum = pgEnum("user_role", ["USER", "ADMIN"]);

export const users = pgTable("users", {
  id: uuid("id").primaryKey().defaultRandom(),

  name: text("name").notNull(),
  email: text("email").notNull().unique(),
  passwordHash: text("password_hash").notNull(),

  role: userRoleEnum("role").notNull().default("USER"),

  isActive: boolean("is_active").notNull().default(true),

  createdAt: timestamp("created_at").notNull().defaultNow(),
  updatedAt: timestamp("updated_at").notNull().defaultNow(),
});

export const refreshTokens = pgTable(
  "refresh_tokens",
  {
    id: uuid("id").primaryKey().defaultRandom(),

    userId: uuid("user_id")
      .notNull()
      .references(() => users.id, { onDelete: "cascade" }),

    tokenHash: text("token_hash").notNull(),

    expiresAt: timestamp("expires_at").notNull(),
    isRevoked: boolean("is_revoked").notNull().default(false),

    createdAt: timestamp("created_at").notNull().defaultNow(),
  },
  (table) => ({
    userIdIdx: index("idx_refresh_tokens_user_id").on(table.userId),
    tokenHashIdx: uniqueIndex("idx_refresh_tokens_token_hash").on(table.tokenHash),
  })
);

export const shortLinks = pgTable(
  "short_links",
  {
    id: uuid("id").primaryKey().defaultRandom(),

    userId: uuid("user_id")
      .notNull()
      .references(() => users.id, { onDelete: "cascade" }),

    originalUrl: text("original_url").notNull(),
    shortCode: text("short_code").notNull().unique(),

    title: text("title"),

    totalClicks: bigint("total_clicks", { mode: "number" })
      .notNull()
      .default(0),

    isActive: boolean("is_active").notNull().default(true),
    expiresAt: timestamp("expires_at"),

    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => ({
    userIdIdx: index("idx_short_links_user_id").on(table.userId),
    shortCodeIdx: uniqueIndex("idx_short_links_short_code").on(table.shortCode),
  })
);

export const clickEvents = pgTable(
  "click_events",
  {
    id: uuid("id").primaryKey().defaultRandom(),

    shortLinkId: uuid("short_link_id")
      .notNull()
      .references(() => shortLinks.id, { onDelete: "cascade" }),

    ipAddress: text("ip_address"),
    userAgent: text("user_agent"),
    referrer: text("referrer"),

    clickedAt: timestamp("clicked_at").notNull().defaultNow(),
  },
  (table) => ({
    shortLinkIdIdx: index("idx_click_events_short_link_id").on(table.shortLinkId),
    clickedAtIdx: index("idx_click_events_clicked_at").on(table.clickedAt),
  })
);

export const idRanges = pgTable(
  "id_ranges",
  {
    id: uuid("id").primaryKey().defaultRandom(),

    serverId: text("server_id").notNull().unique(),

    startId: bigint("start_id", { mode: "number" }).notNull(),
    endId: bigint("end_id", { mode: "number" }).notNull(),
    nextId: bigint("next_id", { mode: "number" }).notNull(),

    isActive: boolean("is_active").notNull().default(true),

    createdAt: timestamp("created_at").notNull().defaultNow(),
  }
);
