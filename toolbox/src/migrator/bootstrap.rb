# frozen_string_literal: true

module Toolbox
  module Migrator
    class Bootstrap
      def initialize(connection:, logger: Logger.new(STDOUT))
        @connection = connection
        @logger = logger
      end

      def call
        @logger.info('Creating migrations table')
        @connection.exec <<~SQL
          BEGIN;

          CREATE TABLE IF NOT EXISTS public.migrations (
            name       text                        PRIMARY KEY,
            created_at timestamp without time zone NOT NULL DEFAULT current_timestamp
          );

          COMMIT;
        SQL
      end
    end
  end
end
