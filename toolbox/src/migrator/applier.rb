# frozen_string_literal: true

module Toolbox
  module Migrator
    class Applier
      def initialize(connection:, logger:)
        @connection = connection
        @logger = logger
      end

      def call(migrations:)
        migrations.each do |migration|
          @logger.info("Applying #{migration[:name]}")

          @connection.exec <<~SQL
            BEGIN;

            #{migration[:sql_up]}

            INSERT INTO migrations (name) VALUES ('#{migration[:name]}');

            COMMIT;
          SQL
        end
      end
    end
  end
end
