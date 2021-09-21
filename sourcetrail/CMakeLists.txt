cmake_minimum_required (VERSION 2.6)
message(STATUS "The SourcetrailDB Go bindings")
# --- Default Flags ---
set(PROJECT_NAME "SourcetrailGoIndexer")
project(${PROJECT_NAME})
set(CMAKE_BUILD_TYPE_INIT "Debug")

if(UNIX)
	set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -pthread -fPIC")
	set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -pthread -fPIC")
endif()

# --- SourcetrailDB Path ---

set(SOURCE_DIR "/home/kaiz/Sourcetrail/SourcetrailDB") 


# --- Version ---
file (STRINGS "${SOURCE_DIR}/version.txt" VERSION_STRING)

string(REGEX REPLACE "v([0-9]+)\\..*" "\\1" INTERFACE_VERSION "${VERSION_STRING}")
string(REGEX REPLACE "v[0-9]+\\.db([0-9]+)\\..*" "\\1" DATABASE_VERSION "${VERSION_STRING}")
string(REGEX REPLACE "v[0-9]+\\.db[0-9]+\\.p([0-9]+).*" "\\1" PATCH_VERSION "${VERSION_STRING}")
set(VERSION_STRING_UNDERSCORE "v${INTERFACE_VERSION}_db${DATABASE_VERSION}_p${PATCH_VERSION}")

message(STATUS "SourcetrailDB Version: ${VERSION_STRING}")
message(STATUS "Interface Version: ${INTERFACE_VERSION}")
message(STATUS "Database Version: ${DATABASE_VERSION}")
message(STATUS "Patch Version: ${PATCH_VERSION}")

# --- Core ---
set(LIB_CORE_TARGET_NAME "lib_core")
set(CORE_SOURCE_DIR "${SOURCE_DIR}/core")
set(CORE_BINARY_DIR "${CMAKE_BINARY_DIR}/core") 

add_subdirectory(${CORE_SOURCE_DIR} ${CORE_BINARY_DIR})

# --- Setup Paths ---
set(RESOURCES_SWIG_DIR "${SOURCE_DIR}/resources_swig")
set(GENERATED_SRC_DIR "${CORE_BINARY_DIR}/src")
set(GO_BINDING_OUTPUT_DIR ${CMAKE_CURRENT_BINARY_DIR} PARENT_SCOPE)

set(SWIG_INTERFACE_FILE "${RESOURCES_SWIG_DIR}/interface/sourcetraildb.i")
set(GO_BINDING_TARGET_NAME "bindings_go")

# --- Find Go ---

find_program(GO_EXECUTABLE go REQUIRED)


# --- Find Swig ---

find_package(SWIG REQUIRED)
include(${SWIG_USE_FILE})

message("pppp: ${GO_INCLUDE_DIRS}")

# --- Configure Target ---

set_source_files_properties(${SWIG_INTERFACE_FILE} PROPERTIES CPLUSPLUS ON)
set(CMAKE_SWIG_FLAGS -cgo -intgosize 64)
include_directories(
	"${RESOURCES_SWIG_DIR}/include"
	"${CORE_SOURCE_DIR}/include"
	${GO_INCLUDE_DIRS}
)

swig_add_library(
	${GO_BINDING_TARGET_NAME}
	TYPE SHARED
	LANGUAGE go
	OUTPUT_DIR ${GO_BINDING_OUTPUT_DIR}
	OUTFILE_DIR ${GENERATED_SRC_DIR}
	SOURCES ${SWIG_INTERFACE_FILE} ${RESOURCES_SWIG_DIR}/src/sourcetraildb.cpp
)

set_target_properties(${GO_BINDING_TARGET_NAME}
                      PROPERTIES OUTPUT_NAME libSourcetraildb)

swig_link_libraries(${GO_BINDING_TARGET_NAME} ${GO_LIBRARIES} ${LIB_CORE_TARGET_NAME})

file(COPY ${RESOURCES_SWIG_DIR}/include/sourcetraildb.h
     DESTINATION ${CMAKE_CURRENT_BINARY_DIR})