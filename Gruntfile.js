module.exports = function (grunt) {

    grunt.initConfig({
			// Watches for changes and runs tasks
		// Livereload is setup for the 35729 port by default
		watch: {
            all: {
    			files: [
                    'css/*.css',
                    'js/*.js',
                    'template/**/*.html'
                ],
				options: {
					livereload: 35732
				}
            }
        }
    });

    grunt.registerTask('default', ['watch:all']);

    grunt.loadNpmTasks('grunt-contrib-watch');
};
