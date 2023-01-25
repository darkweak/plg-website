var gulp = require('gulp'),
    postcss = require('gulp-postcss'),
    concatCss = require('gulp-concat-css'),
    cssnano = require('gulp-cssnano');
    babel = require('gulp-babel');
    minifyjs = require('gulp-js-minify');

gulp.task('dev-css', function () {
    return gulp.src('./css/index.css')
      .pipe(postcss([
        require('tailwindcss'),
        require('autoprefixer'),
      ]))
      .pipe(concatCss('index.css'))
      .pipe(cssnano({
        reduceIdents: false,
        discardComments: {removeAll: true}
      }))
      .pipe(gulp.dest('./static/css/'))
});

gulp.task('build-css', function () {
    return gulp.src('./css/index.css')
      .pipe(postcss([
        require('tailwindcss'),
        require('autoprefixer'),
      ]))
      .pipe(concatCss('index.css'))
      .pipe(cssnano({
        reduceIdents: false,
        discardComments: {removeAll: true}
      }))
      .pipe(gulp.dest('./static/css/'))
});

gulp.task('js', function () {
  return gulp.src("js/*.js")
    .pipe(babel())
    .pipe(minifyjs())
    .pipe(gulp.dest('./static/js/'));
});

gulp.task('watch-css', function() {
  gulp.watch('./css/*.css', gulp.series('dev-css'));
});

gulp.task('dev', gulp.series('dev-css', 'js'));

gulp.task('build', gulp.series('build-css', 'js'));