// Pattern 26: Async sequence patterns
component AsyncController {
    depends on AsyncService : IAsync as asyncSvc
    depends on CallbackHandler : ICallback as callback
    depends on EventBus : IEventBus as events

    provides AsyncAPI {
        async startProcess(id: string) -> ProcessHandle
        async waitForCompletion(handle: ProcessHandle) -> Result
        onComplete(handler: Handler)
    }

    flow AsyncWithAwait {
        await asyncSvc.startLongRunning(taskId)
        handle = asyncSvc.getHandle()
        await asyncSvc.checkStatus(handle)
        status = asyncSvc.getStatus()
        if status.isComplete {
            await asyncSvc.fetchResult(handle)
            result = asyncSvc.getResult()
            return result
        } else {
            await asyncSvc.waitMore(handle)
            await asyncSvc.fetchResult(handle)
            result = asyncSvc.getResult()
            return result
        }
    }

    flow AsyncChain {
        await asyncSvc.step1(input)
        step1Result = asyncSvc.getResult()
        await asyncSvc.step2(step1Result)
        step2Result = asyncSvc.getResult()
        await asyncSvc.step3(step2Result)
        step3Result = asyncSvc.getResult()
        await asyncSvc.finalize(step3Result)
        finalResult = asyncSvc.getResult()
        return finalResult
    }

    flow AsyncWithCallback {
        asyncSvc.startAsync(taskId)
        callback.register(taskId, handler)
        events.emit(taskStartedEvent)
        return taskId
    }

    flow ParallelAsync {
        asyncSvc.startTask(task1Id)
        task1 = asyncSvc.getTask()
        asyncSvc.startTask(task2Id)
        task2 = asyncSvc.getTask()
        asyncSvc.startTask(task3Id)
        task3 = asyncSvc.getTask()
        await asyncSvc.waitFor(task1)
        result1 = asyncSvc.getResult()
        await asyncSvc.waitFor(task2)
        result2 = asyncSvc.getResult()
        await asyncSvc.waitFor(task3)
        result3 = asyncSvc.getResult()
        combined = self.combine(result1, result2, result3)
        return combined
    }
}
