# Builder stage
FROM node:22-alpine AS builder

WORKDIR /app

# Create non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy dependency files and install dependencies
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile

# Copy the rest of the application
COPY . .

# Build the SvelteKit app
RUN yarn build

# Prune dev dependencies for production
RUN yarn install --production --ignore-scripts --prefer-offline

# Production stage
FROM node:22-alpine

WORKDIR /app

# Create non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy built app and production dependencies
COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./

# Set permissions
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 3000
ENV NODE_ENV=production

CMD ["node", "build"]
