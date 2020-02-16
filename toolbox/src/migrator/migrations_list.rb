# frozen_string_literal: true

module Toolbox
  module Migrator
    class MigrationsList
      APPLIED_MIGRATIONS_SQL = <<~SQL
        SELECT * FROM public.migrations ORDER BY name ASC;
      SQL

      def initialize(connection:, logger: Logger.new(STDOUT))
        @connection = connection
        @logger = logger
      end

      def pending(path:)
        migration_files = files(path: path).sort
        latest_migration_applied = applied.last

        migration_files =
          if latest_migration_applied
            latest_migration_time =
              Time.parse(latest_migration_applied[:name].split('-')[0])

            migration_files.filter do |path|
              migration_file_name = File.basename(path)
              migration_file_time = Time.parse(migration_file_name.split('-')[0])

              migration_file_time > latest_migration_time
            end
          else
            migration_files
          end

        migration_files.map do |path|
          { name: File.basename(path) }.tap do |migration|
            migration[:sql_up], migration[:sql_down] =
              File.read(path).split('--migrator:down')
          end
        end
      end

      def applied
        [].tap do |migrations|
          @connection.exec(APPLIED_MIGRATIONS_SQL) do |result|
            result.each do |row|
              migrations << {
                name: row['name'],
                created_at: row['created_at']
              }
            end
          end
        end
      end

      def files(path:)
        Dir["#{path}/*.sql"]
      end
    end
  end
end
