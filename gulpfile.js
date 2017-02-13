"use strict";

const gulp = require("gulp");
const jsHint = require("gulp-jshint");
const babel = require("gulp-babel");
const clean = require("gulp-clean");

gulp.task("lint", function() {
    return gulp.src(["./src/**/*.js"])
        .pipe(jsHint())
        .pipe(jsHint.reporter("default"))
        .pipe(jsHint.reporter("fail")); // fail the task if we fail the linting
});

gulp.task("transpile", function() {
    return gulp.src("./src/**/*.js")
        .pipe(babel())
        .pipe(gulp.dest("target"));
});

gulp.task("clean", function() {
    return gulp.src("./target/**/*", { read: false })
        .pipe(clean());
});

gulp.task("default", ["lint", "transpile"]);
