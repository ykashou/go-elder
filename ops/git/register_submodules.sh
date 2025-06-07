#!/bin/bash
# Add Git submodules with repository URLs
# Adding go-magefile submodule
# Syntax: git submodule add [repository-url] [local-path]
git submodule add https://github.com/audiomage-dev/go-magefile.git pkg/go-magefile

# Adding go-tensor submodule
git submodule add https://github.com/audiomage-dev/go-tensor.git pkg/go-tensor

# Adding go-kernels submodule
git submodule add https://github.com/audiomage-dev/go-kernels.git pkg/go-kernels

# # Adding go-atree submodule
# git submodule add https://github.com/audiomage-dev/go-atree.git pkg/go-atree

# # Adding go-model submodule
# git submodule add https://github.com/audiomage-dev/go-model.git pkg/go-model

# # Adding go-autograd submodule
# git submodule add https://github.com/audiomage-dev/go-autograd.git pkg/go-autograd

# # Adding go-trainer submodule
# git submodule add https://github.com/audiomage-dev/go-trainer.git pkg/go-trainer

# # Adding go-utilities submodule
# git submodule add https://github.com/audiomage-dev/go-utilities.git pkg/go-utilities

# Adding go-benchmarks submodule
git submodule add https://github.com/audiomage-dev/go-benchmarks.git internal/go-benchmarks

# # Adding go-experiments submodule
# git submodule add https://github.com/audiomage-dev/go-experiments.git internal/go-experiments