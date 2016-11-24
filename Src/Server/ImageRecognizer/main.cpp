//
// Created by bentoo on 24.11.16.
//

#include <tensorflow/c_api.h>
#include <cstdio>

int main()
{
	TF_Graph* graph = TF_NewGraph();

	if(graph == nullptr)
		return 1;

	TF_Status* status = TF_NewStatus();
	TF_Session* session = TF_NewSession(TF_NewSessionOptions(), status);

	auto code = TF_GetCode(status);
	if(code != TF_OK)
	{
		printf("%d : %s", code, TF_Message(status));
	}

	return 0;
}
