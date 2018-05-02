// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: storage/engine/enginepb/file_registry.proto

#ifndef PROTOBUF_storage_2fengine_2fenginepb_2ffile_5fregistry_2eproto__INCLUDED
#define PROTOBUF_storage_2fengine_2fenginepb_2ffile_5fregistry_2eproto__INCLUDED

#include <string>

#include <google/protobuf/stubs/common.h>

#if GOOGLE_PROTOBUF_VERSION < 3005000
#error This file was generated by a newer version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please update
#error your headers.
#endif
#if 3005001 < GOOGLE_PROTOBUF_MIN_PROTOC_VERSION
#error This file was generated by an older version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please
#error regenerate this file with a newer version of protoc.
#endif

#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/arena.h>
#include <google/protobuf/arenastring.h>
#include <google/protobuf/generated_message_table_driven.h>
#include <google/protobuf/generated_message_util.h>
#include <google/protobuf/metadata_lite.h>
#include <google/protobuf/message_lite.h>
#include <google/protobuf/repeated_field.h>  // IWYU pragma: export
#include <google/protobuf/extension_set.h>  // IWYU pragma: export
#include <google/protobuf/map.h>  // IWYU pragma: export
#include <google/protobuf/map_entry_lite.h>
#include <google/protobuf/map_field_lite.h>
#include <google/protobuf/generated_enum_util.h>
// @@protoc_insertion_point(includes)

namespace protobuf_storage_2fengine_2fenginepb_2ffile_5fregistry_2eproto {
// Internal implementation detail -- do not use these members.
struct TableStruct {
  static const ::google::protobuf::internal::ParseTableField entries[];
  static const ::google::protobuf::internal::AuxillaryParseTableField aux[];
  static const ::google::protobuf::internal::ParseTable schema[3];
  static const ::google::protobuf::internal::FieldMetadata field_metadata[];
  static const ::google::protobuf::internal::SerializationTable serialization_table[];
  static const ::google::protobuf::uint32 offsets[];
};
void InitDefaultsFileRegistry_FilesEntry_DoNotUseImpl();
void InitDefaultsFileRegistry_FilesEntry_DoNotUse();
void InitDefaultsFileRegistryImpl();
void InitDefaultsFileRegistry();
void InitDefaultsFileEntryImpl();
void InitDefaultsFileEntry();
inline void InitDefaults() {
  InitDefaultsFileRegistry_FilesEntry_DoNotUse();
  InitDefaultsFileRegistry();
  InitDefaultsFileEntry();
}
}  // namespace protobuf_storage_2fengine_2fenginepb_2ffile_5fregistry_2eproto
namespace cockroach {
namespace storage {
namespace engine {
namespace enginepb {
class FileEntry;
class FileEntryDefaultTypeInternal;
extern FileEntryDefaultTypeInternal _FileEntry_default_instance_;
class FileRegistry;
class FileRegistryDefaultTypeInternal;
extern FileRegistryDefaultTypeInternal _FileRegistry_default_instance_;
class FileRegistry_FilesEntry_DoNotUse;
class FileRegistry_FilesEntry_DoNotUseDefaultTypeInternal;
extern FileRegistry_FilesEntry_DoNotUseDefaultTypeInternal _FileRegistry_FilesEntry_DoNotUse_default_instance_;
}  // namespace enginepb
}  // namespace engine
}  // namespace storage
}  // namespace cockroach
namespace cockroach {
namespace storage {
namespace engine {
namespace enginepb {

enum RegistryVersion {
  Base = 0,
  RegistryVersion_INT_MIN_SENTINEL_DO_NOT_USE_ = ::google::protobuf::kint32min,
  RegistryVersion_INT_MAX_SENTINEL_DO_NOT_USE_ = ::google::protobuf::kint32max
};
bool RegistryVersion_IsValid(int value);
const RegistryVersion RegistryVersion_MIN = Base;
const RegistryVersion RegistryVersion_MAX = Base;
const int RegistryVersion_ARRAYSIZE = RegistryVersion_MAX + 1;

enum EnvType {
  Plaintext = 0,
  Store = 1,
  Data = 2,
  EnvType_INT_MIN_SENTINEL_DO_NOT_USE_ = ::google::protobuf::kint32min,
  EnvType_INT_MAX_SENTINEL_DO_NOT_USE_ = ::google::protobuf::kint32max
};
bool EnvType_IsValid(int value);
const EnvType EnvType_MIN = Plaintext;
const EnvType EnvType_MAX = Data;
const int EnvType_ARRAYSIZE = EnvType_MAX + 1;

// ===================================================================

class FileRegistry_FilesEntry_DoNotUse : public ::google::protobuf::internal::MapEntryLite<FileRegistry_FilesEntry_DoNotUse, 
    ::std::string, ::cockroach::storage::engine::enginepb::FileEntry,
    ::google::protobuf::internal::WireFormatLite::TYPE_STRING,
    ::google::protobuf::internal::WireFormatLite::TYPE_MESSAGE,
    0 > {
public:
  typedef ::google::protobuf::internal::MapEntryLite<FileRegistry_FilesEntry_DoNotUse, 
    ::std::string, ::cockroach::storage::engine::enginepb::FileEntry,
    ::google::protobuf::internal::WireFormatLite::TYPE_STRING,
    ::google::protobuf::internal::WireFormatLite::TYPE_MESSAGE,
    0 > SuperType;
  FileRegistry_FilesEntry_DoNotUse();
  FileRegistry_FilesEntry_DoNotUse(::google::protobuf::Arena* arena);
  void MergeFrom(const FileRegistry_FilesEntry_DoNotUse& other);
  static const FileRegistry_FilesEntry_DoNotUse* internal_default_instance() { return reinterpret_cast<const FileRegistry_FilesEntry_DoNotUse*>(&_FileRegistry_FilesEntry_DoNotUse_default_instance_); }
};

// -------------------------------------------------------------------

class FileRegistry : public ::google::protobuf::MessageLite /* @@protoc_insertion_point(class_definition:cockroach.storage.engine.enginepb.FileRegistry) */ {
 public:
  FileRegistry();
  virtual ~FileRegistry();

  FileRegistry(const FileRegistry& from);

  inline FileRegistry& operator=(const FileRegistry& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  FileRegistry(FileRegistry&& from) noexcept
    : FileRegistry() {
    *this = ::std::move(from);
  }

  inline FileRegistry& operator=(FileRegistry&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const FileRegistry& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const FileRegistry* internal_default_instance() {
    return reinterpret_cast<const FileRegistry*>(
               &_FileRegistry_default_instance_);
  }
  static PROTOBUF_CONSTEXPR int const kIndexInFileMessages =
    1;

  void Swap(FileRegistry* other);
  friend void swap(FileRegistry& a, FileRegistry& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline FileRegistry* New() const PROTOBUF_FINAL { return New(NULL); }

  FileRegistry* New(::google::protobuf::Arena* arena) const PROTOBUF_FINAL;
  void CheckTypeAndMergeFrom(const ::google::protobuf::MessageLite& from)
    PROTOBUF_FINAL;
  void CopyFrom(const FileRegistry& from);
  void MergeFrom(const FileRegistry& from);
  void Clear() PROTOBUF_FINAL;
  bool IsInitialized() const PROTOBUF_FINAL;

  size_t ByteSizeLong() const PROTOBUF_FINAL;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) PROTOBUF_FINAL;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const PROTOBUF_FINAL;
  void DiscardUnknownFields();
  int GetCachedSize() const PROTOBUF_FINAL { return _cached_size_; }
  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const;
  void InternalSwap(FileRegistry* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::std::string GetTypeName() const PROTOBUF_FINAL;

  // nested types ----------------------------------------------------


  // accessors -------------------------------------------------------

  // map<string, .cockroach.storage.engine.enginepb.FileEntry> files = 2;
  int files_size() const;
  void clear_files();
  static const int kFilesFieldNumber = 2;
  const ::google::protobuf::Map< ::std::string, ::cockroach::storage::engine::enginepb::FileEntry >&
      files() const;
  ::google::protobuf::Map< ::std::string, ::cockroach::storage::engine::enginepb::FileEntry >*
      mutable_files();

  // .cockroach.storage.engine.enginepb.RegistryVersion version = 1;
  void clear_version();
  static const int kVersionFieldNumber = 1;
  ::cockroach::storage::engine::enginepb::RegistryVersion version() const;
  void set_version(::cockroach::storage::engine::enginepb::RegistryVersion value);

  // @@protoc_insertion_point(class_scope:cockroach.storage.engine.enginepb.FileRegistry)
 private:

  ::google::protobuf::internal::InternalMetadataWithArenaLite _internal_metadata_;
  ::google::protobuf::internal::MapFieldLite<
      FileRegistry_FilesEntry_DoNotUse,
      ::std::string, ::cockroach::storage::engine::enginepb::FileEntry,
      ::google::protobuf::internal::WireFormatLite::TYPE_STRING,
      ::google::protobuf::internal::WireFormatLite::TYPE_MESSAGE,
      0 > files_;
  int version_;
  mutable int _cached_size_;
  friend struct ::protobuf_storage_2fengine_2fenginepb_2ffile_5fregistry_2eproto::TableStruct;
  friend void ::protobuf_storage_2fengine_2fenginepb_2ffile_5fregistry_2eproto::InitDefaultsFileRegistryImpl();
};
// -------------------------------------------------------------------

class FileEntry : public ::google::protobuf::MessageLite /* @@protoc_insertion_point(class_definition:cockroach.storage.engine.enginepb.FileEntry) */ {
 public:
  FileEntry();
  virtual ~FileEntry();

  FileEntry(const FileEntry& from);

  inline FileEntry& operator=(const FileEntry& from) {
    CopyFrom(from);
    return *this;
  }
  #if LANG_CXX11
  FileEntry(FileEntry&& from) noexcept
    : FileEntry() {
    *this = ::std::move(from);
  }

  inline FileEntry& operator=(FileEntry&& from) noexcept {
    if (GetArenaNoVirtual() == from.GetArenaNoVirtual()) {
      if (this != &from) InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }
  #endif
  static const FileEntry& default_instance();

  static void InitAsDefaultInstance();  // FOR INTERNAL USE ONLY
  static inline const FileEntry* internal_default_instance() {
    return reinterpret_cast<const FileEntry*>(
               &_FileEntry_default_instance_);
  }
  static PROTOBUF_CONSTEXPR int const kIndexInFileMessages =
    2;

  void Swap(FileEntry* other);
  friend void swap(FileEntry& a, FileEntry& b) {
    a.Swap(&b);
  }

  // implements Message ----------------------------------------------

  inline FileEntry* New() const PROTOBUF_FINAL { return New(NULL); }

  FileEntry* New(::google::protobuf::Arena* arena) const PROTOBUF_FINAL;
  void CheckTypeAndMergeFrom(const ::google::protobuf::MessageLite& from)
    PROTOBUF_FINAL;
  void CopyFrom(const FileEntry& from);
  void MergeFrom(const FileEntry& from);
  void Clear() PROTOBUF_FINAL;
  bool IsInitialized() const PROTOBUF_FINAL;

  size_t ByteSizeLong() const PROTOBUF_FINAL;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input) PROTOBUF_FINAL;
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const PROTOBUF_FINAL;
  void DiscardUnknownFields();
  int GetCachedSize() const PROTOBUF_FINAL { return _cached_size_; }
  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const;
  void InternalSwap(FileEntry* other);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return NULL;
  }
  inline void* MaybeArenaPtr() const {
    return NULL;
  }
  public:

  ::std::string GetTypeName() const PROTOBUF_FINAL;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // bytes encryption_settings = 2;
  void clear_encryption_settings();
  static const int kEncryptionSettingsFieldNumber = 2;
  const ::std::string& encryption_settings() const;
  void set_encryption_settings(const ::std::string& value);
  #if LANG_CXX11
  void set_encryption_settings(::std::string&& value);
  #endif
  void set_encryption_settings(const char* value);
  void set_encryption_settings(const void* value, size_t size);
  ::std::string* mutable_encryption_settings();
  ::std::string* release_encryption_settings();
  void set_allocated_encryption_settings(::std::string* encryption_settings);

  // .cockroach.storage.engine.enginepb.EnvType env_type = 1;
  void clear_env_type();
  static const int kEnvTypeFieldNumber = 1;
  ::cockroach::storage::engine::enginepb::EnvType env_type() const;
  void set_env_type(::cockroach::storage::engine::enginepb::EnvType value);

  // @@protoc_insertion_point(class_scope:cockroach.storage.engine.enginepb.FileEntry)
 private:

  ::google::protobuf::internal::InternalMetadataWithArenaLite _internal_metadata_;
  ::google::protobuf::internal::ArenaStringPtr encryption_settings_;
  int env_type_;
  mutable int _cached_size_;
  friend struct ::protobuf_storage_2fengine_2fenginepb_2ffile_5fregistry_2eproto::TableStruct;
  friend void ::protobuf_storage_2fengine_2fenginepb_2ffile_5fregistry_2eproto::InitDefaultsFileEntryImpl();
};
// ===================================================================


// ===================================================================

#ifdef __GNUC__
  #pragma GCC diagnostic push
  #pragma GCC diagnostic ignored "-Wstrict-aliasing"
#endif  // __GNUC__
// -------------------------------------------------------------------

// FileRegistry

// .cockroach.storage.engine.enginepb.RegistryVersion version = 1;
inline void FileRegistry::clear_version() {
  version_ = 0;
}
inline ::cockroach::storage::engine::enginepb::RegistryVersion FileRegistry::version() const {
  // @@protoc_insertion_point(field_get:cockroach.storage.engine.enginepb.FileRegistry.version)
  return static_cast< ::cockroach::storage::engine::enginepb::RegistryVersion >(version_);
}
inline void FileRegistry::set_version(::cockroach::storage::engine::enginepb::RegistryVersion value) {
  
  version_ = value;
  // @@protoc_insertion_point(field_set:cockroach.storage.engine.enginepb.FileRegistry.version)
}

// map<string, .cockroach.storage.engine.enginepb.FileEntry> files = 2;
inline int FileRegistry::files_size() const {
  return files_.size();
}
inline void FileRegistry::clear_files() {
  files_.Clear();
}
inline const ::google::protobuf::Map< ::std::string, ::cockroach::storage::engine::enginepb::FileEntry >&
FileRegistry::files() const {
  // @@protoc_insertion_point(field_map:cockroach.storage.engine.enginepb.FileRegistry.files)
  return files_.GetMap();
}
inline ::google::protobuf::Map< ::std::string, ::cockroach::storage::engine::enginepb::FileEntry >*
FileRegistry::mutable_files() {
  // @@protoc_insertion_point(field_mutable_map:cockroach.storage.engine.enginepb.FileRegistry.files)
  return files_.MutableMap();
}

// -------------------------------------------------------------------

// FileEntry

// .cockroach.storage.engine.enginepb.EnvType env_type = 1;
inline void FileEntry::clear_env_type() {
  env_type_ = 0;
}
inline ::cockroach::storage::engine::enginepb::EnvType FileEntry::env_type() const {
  // @@protoc_insertion_point(field_get:cockroach.storage.engine.enginepb.FileEntry.env_type)
  return static_cast< ::cockroach::storage::engine::enginepb::EnvType >(env_type_);
}
inline void FileEntry::set_env_type(::cockroach::storage::engine::enginepb::EnvType value) {
  
  env_type_ = value;
  // @@protoc_insertion_point(field_set:cockroach.storage.engine.enginepb.FileEntry.env_type)
}

// bytes encryption_settings = 2;
inline void FileEntry::clear_encryption_settings() {
  encryption_settings_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& FileEntry::encryption_settings() const {
  // @@protoc_insertion_point(field_get:cockroach.storage.engine.enginepb.FileEntry.encryption_settings)
  return encryption_settings_.GetNoArena();
}
inline void FileEntry::set_encryption_settings(const ::std::string& value) {
  
  encryption_settings_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:cockroach.storage.engine.enginepb.FileEntry.encryption_settings)
}
#if LANG_CXX11
inline void FileEntry::set_encryption_settings(::std::string&& value) {
  
  encryption_settings_.SetNoArena(
    &::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::move(value));
  // @@protoc_insertion_point(field_set_rvalue:cockroach.storage.engine.enginepb.FileEntry.encryption_settings)
}
#endif
inline void FileEntry::set_encryption_settings(const char* value) {
  GOOGLE_DCHECK(value != NULL);
  
  encryption_settings_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:cockroach.storage.engine.enginepb.FileEntry.encryption_settings)
}
inline void FileEntry::set_encryption_settings(const void* value, size_t size) {
  
  encryption_settings_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:cockroach.storage.engine.enginepb.FileEntry.encryption_settings)
}
inline ::std::string* FileEntry::mutable_encryption_settings() {
  
  // @@protoc_insertion_point(field_mutable:cockroach.storage.engine.enginepb.FileEntry.encryption_settings)
  return encryption_settings_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* FileEntry::release_encryption_settings() {
  // @@protoc_insertion_point(field_release:cockroach.storage.engine.enginepb.FileEntry.encryption_settings)
  
  return encryption_settings_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void FileEntry::set_allocated_encryption_settings(::std::string* encryption_settings) {
  if (encryption_settings != NULL) {
    
  } else {
    
  }
  encryption_settings_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), encryption_settings);
  // @@protoc_insertion_point(field_set_allocated:cockroach.storage.engine.enginepb.FileEntry.encryption_settings)
}

#ifdef __GNUC__
  #pragma GCC diagnostic pop
#endif  // __GNUC__
// -------------------------------------------------------------------

// -------------------------------------------------------------------


// @@protoc_insertion_point(namespace_scope)

}  // namespace enginepb
}  // namespace engine
}  // namespace storage
}  // namespace cockroach

namespace google {
namespace protobuf {

template <> struct is_proto_enum< ::cockroach::storage::engine::enginepb::RegistryVersion> : ::google::protobuf::internal::true_type {};
template <> struct is_proto_enum< ::cockroach::storage::engine::enginepb::EnvType> : ::google::protobuf::internal::true_type {};

}  // namespace protobuf
}  // namespace google

// @@protoc_insertion_point(global_scope)

#endif  // PROTOBUF_storage_2fengine_2fenginepb_2ffile_5fregistry_2eproto__INCLUDED
