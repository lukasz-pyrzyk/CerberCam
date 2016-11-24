# Find Tensorflow
#
# TENSORFLOW_INCLUDE_DIR
# TENSORFLOW_LIBRARY
# TENSORFLOW_FOUND

find_path(TENSORFLOW_INCLUDE_DIR NAMES tensorflow/c_api.h HINTS
        "$ENV{TENSORFLOW_PATH}"
        "$ENV{TENSORFLOW_PATH}/include")
find_library(TENSORFLOW_LIBRARY NAMES tensorflow HINTS
        "$ENV{TENSORFLOW_PATH}/lib")

include(FindPackageHandleStandardArgs)

find_package_handle_standard_args(Tensorflow DEFAULT_MSG TENSORFLOW_LIBRARY TENSORFLOW_INCLUDE_DIR)

mark_as_advanced(TENSORFLOW_INCLUDE_DIR TENSORFLOW_LIBRARY)