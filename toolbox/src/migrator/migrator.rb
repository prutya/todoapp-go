# frozen_string_literal: true

require 'logger'
require 'time'

require 'pg'

module Toolbox
  module Migrator
    class Config
      attr_reader :db_url
      attr_reader :migrations_dir

      def initialize(env: ENV)
        @db_url = env['TOOLBOX_DB_URL']
        @migrations_dir = File.expand_path('../../migrations', __dir__)
      end
    end
  end
end

require_relative './bootstrap'
require_relative './migrations_list'
require_relative './applier'
