#include "stdio.h"
#include "string.h"
#include "stdlib.h"
#include "windows.h"
#include "include/dart_api.h"
#include "include/dart_native_api.h"

typedef  unsigned long DWORD;
Dart_NativeFunction ResolveName(Dart_Handle name,
								int argc,
								bool* auto_setup_scope);


DART_EXPORT Dart_Handle mystdin_Init(Dart_Handle parent_library) {
	if (Dart_IsError(parent_library)) {
		return parent_library;
	}

	Dart_Handle result_code =
		Dart_SetNativeResolver(parent_library, ResolveName, NULL);
	if (Dart_IsError(result_code)) {
		return result_code;
	}

	return Dart_Null();
}


Dart_Handle HandleError(Dart_Handle handle) {
	if (Dart_IsError(handle)) {
		Dart_PropagateError(handle);
	}
	return handle;
}

uint8_t* readBytes(DWORD* read) {
	HANDLE h = GetStdHandle(STD_INPUT_HANDLE);
	const int size = 64 * 1024;
	uint8_t* buffer = reinterpret_cast<uint8_t*>(malloc(size));
	bool success = ReadFile(h, buffer, size, read, NULL);
	if (!success && (GetLastError() != ERROR_BROKEN_PIPE)) {
		return NULL;
	}
	DWORD length = *read;
	if (size != length){
		uint8_t* end = reinterpret_cast<uint8_t*>(malloc(length));
		memcpy(end,buffer,sizeof(uint8_t)*length);
		free(buffer);
		return end;
	}
	return buffer;
}

uint8_t* ReadBytes(DWORD* read){
	DWORD size = 0;
	uint8_t *buffer,*buf,*newbuf;
	buf = readBytes(read);
	size = *read;
	newbuf = reinterpret_cast<uint8_t*>(malloc(*read));
	memcpy(newbuf,buf,sizeof(uint8_t)*(*read));
	free(buf);
	while(*read == 64 * 1024){
		buf = readBytes(read);
		buffer = newbuf;
		newbuf = reinterpret_cast<uint8_t*>(malloc(size + *read));
		memcpy(newbuf,buffer,sizeof(uint8_t)*size);
		memcpy(newbuf + size*sizeof(uint8_t),buf,sizeof(uint8_t)*(*read));
		size += *read;
		free(buffer);
		free(buf);
	}
	*read = size;
	return newbuf;
}

void wrappedReadBytes(Dart_Port dest_port_id,
					  Dart_CObject* message) {
						  Dart_Port reply_port_id = ILLEGAL_PORT;
						  if (message->type == Dart_CObject_kArray &&
							  1 == message->value.as_array.length) {
								  // Use .as_array and .as_int32 to access the data in the Dart_CObject.
								  Dart_CObject* param0 = message->value.as_array.values[0];
								  if (param0->type == Dart_CObject_kSendPort) {
									  reply_port_id = param0->value.as_send_port.id;
									  DWORD read = 0;
									  uint8_t* buffer = ReadBytes(&read);

									  if (buffer != NULL) {
										  Dart_CObject result;
										  result.type = Dart_CObject_kTypedData;
										  result.value.as_typed_data.type = Dart_TypedData_kUint8;
										  result.value.as_typed_data.values = buffer;
										  result.value.as_typed_data.length = read;
										  Dart_PostCObject(reply_port_id, &result);
										  free(buffer);
										  // It is OK that result is destroyed when function exits.
										  // Dart_PostCObject has copied its data.
										  return;
									  }
								  }
						  }
						  Dart_CObject result;
						  result.type = Dart_CObject_kNull;
						  Dart_PostCObject(reply_port_id, &result);
}


void ReadBytesServicePort(Dart_NativeArguments arguments) {
	Dart_EnterScope();
	Dart_SetReturnValue(arguments, Dart_Null());
	Dart_Port service_port =
		Dart_NewNativePort("ReadBytesService", wrappedReadBytes, true);
	if (service_port != ILLEGAL_PORT) {
		Dart_Handle send_port = HandleError(Dart_NewSendPort(service_port));
		Dart_SetReturnValue(arguments, send_port);
	}
	Dart_ExitScope();
}


struct FunctionLookup {
	const char* name;
	Dart_NativeFunction function;
};


FunctionLookup function_list[] = {
	{"ReadBytes_ServicePort", ReadBytesServicePort},
	{NULL, NULL}};

	Dart_NativeFunction ResolveName(Dart_Handle name,
		int argc,
		bool* auto_setup_scope) {
			if (!Dart_IsString(name)) {
				return NULL;
			}
			Dart_NativeFunction result = NULL;
			if (auto_setup_scope == NULL) {
				return NULL;
			}

			Dart_EnterScope();
			const char* cname;
			HandleError(Dart_StringToCString(name, &cname));

			for (int i=0; function_list[i].name != NULL; ++i) {
				if (strcmp(function_list[i].name, cname) == 0) {
					*auto_setup_scope = true;
					result = function_list[i].function;
					break;
				}
			}

			if (result != NULL) {
				Dart_ExitScope();
				return result;
			}

			Dart_ExitScope();
			return result;
	}
